package service

import (
	"github.com/wyw14/cry-161/internal/model"
	"testing"
	"time"
)

func validSample(runID model.ID) model.ThicknessSample {
	return model.ThicknessSample{RunID: runID, Frequency: 10, Transmission: 0.95, Nanometers: 42, Valid: true, At: time.Now().UTC()}
}

func TestCaptureInvalidOnlyRaisesIncident(t *testing.T) {
	runtime := NewRuntime()
	_, ch := runtime.Broadcast.Subscribe(1)
	runID := model.NewID()
	invalid := validSample(runID)
	invalid.Valid = false
	if _, err := runtime.Capture(invalid, 1); err == nil {
		t.Fatal("expected validation error")
	}
	if got := len(runtime.Monitor.Samples()); got != 0 {
		t.Fatalf("monitor polluted: %d samples", got)
	}
	if got := len(runtime.Window.Values()); got != 0 {
		t.Fatalf("window polluted: %d samples", got)
	}
	if got := len(runtime.Traces.Trace(runID).Thickness); got != 0 {
		t.Fatalf("traces polluted: %d samples", got)
	}
	if got := runtime.Timeline.Count(); got != 0 {
		t.Fatalf("timeline polluted: %d entries", got)
	}
	select {
	case sample := <-ch:
		t.Fatalf("broadcast polluted: %+v", sample)
	default:
	}
	if got := len(runtime.Incidents.ForRun(runID)); got != 1 {
		t.Fatalf("expected 1 incident, got %d", got)
	}
}

func TestCaptureCommitsValidOnce(t *testing.T) {
	runtime := NewRuntime()
	_, ch := runtime.Broadcast.Subscribe(4)
	runID := model.NewID()
	sample := validSample(runID)
	first, err := runtime.Capture(sample, 1)
	if err != nil {
		t.Fatalf("valid capture rejected: %v", err)
	}
	second, err := runtime.Capture(sample, 1)
	if err != nil {
		t.Fatalf("duplicate capture errored: %v", err)
	}
	if second != first {
		t.Fatalf("duplicate returned different evidence: %+v vs %+v", second, first)
	}
	if got := len(runtime.Monitor.Samples()); got != 1 {
		t.Fatalf("monitor wrote %d samples", got)
	}
	if got := len(runtime.Window.Values()); got != 1 {
		t.Fatalf("window wrote %d samples", got)
	}
	if got := len(runtime.Traces.Trace(runID).Thickness); got != 1 {
		t.Fatalf("traces wrote %d samples", got)
	}
	if got := runtime.Timeline.Count(); got != 1 {
		t.Fatalf("timeline wrote %d entries", got)
	}
	for i := 0; i < 2; i++ {
		select {
		case <-ch:
			if i == 1 {
				t.Fatal("duplicate capture re-broadcast")
			}
		default:
			if i == 0 {
				t.Fatal("valid capture not broadcast")
			}
		}
	}
	if got := len(runtime.Incidents.ForRun(runID)); got != 0 {
		t.Fatalf("unexpected incidents: %d", got)
	}
}

func TestCaptureRetryAfterInvalidCommits(t *testing.T) {
	runtime := NewRuntime()
	runID := model.NewID()
	invalid := validSample(runID)
	invalid.Valid = false
	if _, err := runtime.Capture(invalid, 1); err == nil {
		t.Fatal("expected validation error")
	}
	if _, err := runtime.Capture(validSample(runID), 1); err != nil {
		t.Fatalf("valid retry rejected: %v", err)
	}
	if got := len(runtime.Traces.Trace(runID).Thickness); got != 1 {
		t.Fatalf("expected 1 committed sample, got %d", got)
	}
}
