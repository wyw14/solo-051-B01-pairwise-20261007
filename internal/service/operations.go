package service

import (
	"context"
	"fmt"
	"github.com/wyw14/cry-161/internal/interlock"
	"github.com/wyw14/cry-161/internal/journal"
	"github.com/wyw14/cry-161/internal/model"
	"github.com/wyw14/cry-161/internal/optical"
	"github.com/wyw14/cry-161/internal/quartz"
	"github.com/wyw14/cry-161/internal/recipe"
	"github.com/wyw14/cry-161/internal/run"
	"github.com/wyw14/cry-161/internal/shutter"
	"github.com/wyw14/cry-161/internal/source"
	"github.com/wyw14/cry-161/internal/vacuum"
	"sync"
	"time"
)

type captureKey struct {
	run   model.ID
	layer int
	at    time.Time
}

type captureRecord struct {
	evidence model.LayerEvidence
	errText  string
}

type Runtime struct {
	mu              sync.RWMutex
	Monitor         *optical.Monitor
	Window          *optical.ThicknessWindow
	Broadcast       *optical.Broadcaster
	Traces          *journal.TraceStore
	Incidents       *journal.IncidentStore
	Timeline        *run.Timeline
	VacuumTelemetry *vacuum.Telemetry
	RevisionStore   *recipe.RevisionStore
	InterlockAudit  *interlock.Audit
	Shutter         *shutter.Controller
	Sequence        int64
	StartedAt       time.Time
	captured        map[captureKey]captureRecord
}

func NewRuntime() *Runtime {
	driver := shutter.NewDriver()
	feedback := shutter.NewFeedback()
	return &Runtime{Monitor: optical.NewMonitor(), Window: optical.NewThicknessWindow(128), Broadcast: optical.NewBroadcaster(), Traces: journal.NewTraceStore(), Incidents: journal.NewIncidentStore(), Timeline: run.NewTimeline(), VacuumTelemetry: vacuum.NewTelemetry(128), RevisionStore: recipe.NewRevisionStore(), InterlockAudit: interlock.NewAudit(), Shutter: shutter.NewController(driver, feedback), StartedAt: time.Now().UTC(), captured: map[captureKey]captureRecord{}}
}

func (r *Runtime) NextSequence() int64 {
	r.mu.Lock()
	r.Sequence++
	value := r.Sequence
	r.mu.Unlock()
	return value
}

// Capture validates the sample before touching any projection: an invalid
// sample produces a single incident event and is never written to the
// monitor, window, trace store, timeline, or broadcast to consumers. A valid
// sample is committed to every projection first and only broadcast after the
// commit succeeds. Repeating Capture with the same sample is a no-op that
// replays the original result without writing anything twice.
func (r *Runtime) Capture(sample model.ThicknessSample, layer int) (model.LayerEvidence, error) {
	if sample.At.IsZero() {
		sample.At = time.Now().UTC()
	}
	evidence := model.BuildLayerEvidence(sample, layer)
	validateErr := evidence.Validate()
	key := captureKey{run: sample.RunID, layer: layer, at: sample.At.UTC()}

	r.mu.Lock()
	if prior, ok := r.captured[key]; ok {
		r.mu.Unlock()
		var err error
		if prior.errText != "" {
			err = fmt.Errorf("%s", prior.errText)
		}
		return prior.evidence, err
	}
	if validateErr != nil {
		r.captured[key] = captureRecord{evidence: evidence, errText: validateErr.Error()}
		r.mu.Unlock()
		r.Incidents.Raise(sample.RunID, "warning", validateErr.Error())
		return evidence, validateErr
	}
	r.Monitor.Add(sample)
	r.Window.Add(sample)
	r.Traces.AppendThickness(sample)
	r.Timeline.Add(sample.RunID, "thickness", fmt.Sprintf("%.3f nm", sample.Nanometers))
	r.captured[key] = captureRecord{evidence: evidence}
	r.mu.Unlock()

	r.Broadcast.Publish(sample)
	return evidence, nil
}

func (r *Runtime) RecordFrequency(runID model.ID, values []float64) {
	copied := append([]float64(nil), values...)
	r.Traces.AppendFrequency(runID, copied)
}

func (r *Runtime) Command(runID model.ID, revision int, kind model.CommandKind) (model.Command, error) {
	command := model.NewCommand(runID, revision, r.NextSequence(), kind)
	if err := command.Validate(); err != nil {
		return command.Fail(err), err
	}
	r.Timeline.Add(runID, "command", string(kind))
	return command.Complete("accepted"), nil
}

func (r *Runtime) Summary() map[string]any {
	latest := r.Monitor.Latest()
	return map[string]any{"started_at": r.StartedAt, "sequence": r.Sequence, "latest_thickness": latest, "trace_count": r.Traces.Count(), "timeline_entries": r.Timeline.Count(), "incidents": r.Incidents.List(), "vacuum": r.VacuumTelemetry.Readings(), "shutter_commands": r.Shutter.Commands()}
}

func (s *Service) RunLayer(ctx context.Context, runtime *Runtime, runID model.ID, layer model.Layer) error {
	value, ok := s.Runs.Get(runID)
	if !ok {
		return fmt.Errorf("run not found")
	}
	command, err := runtime.Command(runID, value.Layer, model.CommandCondition)
	if err != nil || command.Status != "completed" {
		return fmt.Errorf("layer command rejected")
	}
	worker := source.NewWorker(s.Source)
	if err := worker.Execute(ctx, layer); err != nil {
		incident := runtime.Incidents.Raise(runID, "critical", err.Error())
		_ = s.Journal.Append(journal.IncidentRaised(incident))
		s.Runs.UpdateError(runID, err)
		return err
	}
	profile := quartz.DefaultCalibration()
	sample := model.ThicknessSample{RunID: runID, Frequency: layer.TargetNanometers / profile.Slope, Transmission: 0.95, Nanometers: layer.TargetNanometers, Valid: true, At: time.Now().UTC()}
	if _, err := runtime.Capture(sample, layer.Number); err != nil {
		return err
	}
	return nil
}
