package service

import (
	"testing"
	"time"

	"github.com/wyw14/cry-161/internal/model"
)

func validSample(runID model.ID, at time.Time) model.ThicknessSample {
	return model.ThicknessSample{
		RunID:        runID,
		Frequency:    1200,
		Transmission: 0.95,
		Nanometers:   100,
		Valid:        true,
		At:           at,
	}
}

func invalidSample(runID model.ID, at time.Time) model.ThicknessSample {
	sample := validSample(runID, at)
	sample.Valid = false
	return sample
}

func TestCaptureInvalidSampleOnlyRaisesEvent(t *testing.T) {
	runtime := NewRuntime()
	runID := model.NewID()
	at := time.Date(2026, 10, 7, 10, 0, 0, 0, time.UTC)

	_, ch := runtime.Broadcast.Subscribe(2)

	if _, err := runtime.Capture(invalidSample(runID, at), 1); err == nil {
		t.Fatal("expected invalid sample to be rejected")
	}

	if len(runtime.Monitor.Samples()) != 0 {
		t.Fatalf("monitor must not store invalid sample, got %d", len(runtime.Monitor.Samples()))
	}
	if !runtime.Monitor.Latest().At.IsZero() {
		t.Fatal("monitor latest must remain zero after invalid sample")
	}
	if len(runtime.Window.Values()) != 0 {
		t.Fatalf("window must not store invalid sample, got %d", len(runtime.Window.Values()))
	}
	if len(runtime.Traces.Trace(runID).Thickness) != 0 {
		t.Fatal("trace history must not contain invalid sample")
	}
	if len(runtime.Timeline.Entries(runID)) != 0 {
		t.Fatal("timeline must not record invalid sample")
	}
	if incidents := runtime.Incidents.ForRun(runID); len(incidents) != 1 {
		t.Fatalf("invalid sample must raise exactly one incident, got %d", len(incidents))
	}
	select {
	case got := <-ch:
		t.Fatalf("invalid sample must not be broadcast, got %+v", got)
	default:
	}
}

func TestCaptureValidSampleCommitsAndBroadcasts(t *testing.T) {
	runtime := NewRuntime()
	runID := model.NewID()
	at := time.Date(2026, 10, 7, 10, 1, 0, 0, time.UTC)

	_, ch := runtime.Broadcast.Subscribe(2)

	evidence, err := runtime.Capture(validSample(runID, at), 1)
	if err != nil {
		t.Fatalf("valid sample must be accepted, got %v", err)
	}
	if evidence.Status != model.EvidenceComplete {
		t.Fatalf("evidence status = %s, want %s", evidence.Status, model.EvidenceComplete)
	}
	if len(runtime.Monitor.Samples()) != 1 {
		t.Fatalf("monitor must hold one sample, got %d", len(runtime.Monitor.Samples()))
	}
	if len(runtime.Window.Values()) != 1 {
		t.Fatalf("window must hold one sample, got %d", len(runtime.Window.Values()))
	}
	if trace := runtime.Traces.Trace(runID); len(trace.Thickness) != 1 {
		t.Fatalf("trace must hold one sample, got %d", len(trace.Thickness))
	}
	if len(runtime.Timeline.Entries(runID)) != 1 {
		t.Fatal("timeline must record one thickness entry")
	}
	select {
	case <-ch:
	default:
		t.Fatal("valid sample must be broadcast after commit")
	}
}

func TestCaptureRepeatedValidSampleWritesOnce(t *testing.T) {
	runtime := NewRuntime()
	runID := model.NewID()
	at := time.Date(2026, 10, 7, 10, 2, 0, 0, time.UTC)
	sample := validSample(runID, at)

	_, ch := runtime.Broadcast.Subscribe(4)

	first, err := runtime.Capture(sample, 1)
	if err != nil {
		t.Fatalf("first capture failed: %v", err)
	}
	second, err := runtime.Capture(sample, 1)
	if err != nil {
		t.Fatalf("repeated capture must replay success, got %v", err)
	}
	if first != second {
		t.Fatal("repeated capture must replay the original evidence")
	}

	if len(runtime.Monitor.Samples()) != 1 {
		t.Fatalf("monitor must hold one sample, got %d", len(runtime.Monitor.Samples()))
	}
	if len(runtime.Window.Values()) != 1 {
		t.Fatalf("window must hold one sample, got %d", len(runtime.Window.Values()))
	}
	if len(runtime.Traces.Trace(runID).Thickness) != 1 {
		t.Fatal("trace must hold one sample")
	}
	if len(runtime.Timeline.Entries(runID)) != 1 {
		t.Fatal("timeline must hold one entry")
	}
	if len(ch) != 1 {
		t.Fatalf("repeated capture must not broadcast twice, channel len = %d", len(ch))
	}
}

func TestCaptureRepeatedInvalidSampleRaisesOnce(t *testing.T) {
	runtime := NewRuntime()
	runID := model.NewID()
	at := time.Date(2026, 10, 7, 10, 3, 0, 0, time.UTC)
	sample := invalidSample(runID, at)

	_, ch := runtime.Broadcast.Subscribe(4)

	if _, err := runtime.Capture(sample, 1); err == nil {
		t.Fatal("first capture should be rejected")
	}
	_, err := runtime.Capture(sample, 1)
	if err == nil {
		t.Fatal("repeated capture must replay the rejection error")
	}

	if incidents := runtime.Incidents.ForRun(runID); len(incidents) != 1 {
		t.Fatalf("repeated invalid capture must not raise again, got %d incidents", len(incidents))
	}
	if len(runtime.Monitor.Samples()) != 0 || len(runtime.Window.Values()) != 0 {
		t.Fatal("invalid repeat must not touch projections")
	}
	if len(ch) != 0 {
		t.Fatalf("invalid repeat must not broadcast, channel len = %d", len(ch))
	}

	// A distinct invalid sample is a new event, not a repeat.
	if _, err := runtime.Capture(invalidSample(runID, at.Add(time.Second)), 1); err == nil {
		t.Fatal("distinct invalid sample should still be rejected")
	}
	if incidents := runtime.Incidents.ForRun(runID); len(incidents) != 2 {
		t.Fatalf("distinct invalid sample should raise a second incident, got %d", len(incidents))
	}
}

func TestCaptureRepeatSameLayerDifferentTimestampCommitsAgain(t *testing.T) {
	runtime := NewRuntime()
	runID := model.NewID()
	at := time.Date(2026, 10, 7, 10, 4, 0, 0, time.UTC)

	if _, err := runtime.Capture(validSample(runID, at), 1); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.Capture(validSample(runID, at.Add(time.Second)), 1); err != nil {
		t.Fatal(err)
	}
	if got := len(runtime.Traces.Trace(runID).Thickness); got != 2 {
		t.Fatalf("distinct timestamps are distinct samples, trace len = %d", got)
	}

	// An invalid repeat replays the same validation error without side effects.
	invalidAt := at.Add(2 * time.Second)
	if _, err := runtime.Capture(invalidSample(runID, invalidAt), 1); err == nil {
		t.Fatal("invalid capture should be rejected")
	}
	if _, err := runtime.Capture(invalidSample(runID, invalidAt), 1); err == nil {
		t.Fatal("repeated invalid capture should replay the rejection error")
	}
	if incidents := runtime.Incidents.ForRun(runID); len(incidents) != 1 {
		t.Fatalf("repeated invalid capture must raise only once, got %d", len(incidents))
	}
}
