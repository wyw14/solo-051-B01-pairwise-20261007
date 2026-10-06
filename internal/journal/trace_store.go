package journal

import (
	"github.com/wyw14/cry-161/internal/model"
	"sync"
)

type Trace struct {
	RunID       model.ID                `json:"run_id"`
	Frequencies []float64               `json:"frequencies"`
	Thickness   []model.ThicknessSample `json:"thickness"`
}

type TraceStore struct {
	mu     sync.RWMutex
	traces map[model.ID]Trace
}

func NewTraceStore() *TraceStore {
	return &TraceStore{traces: map[model.ID]Trace{}}
}

func (s *TraceStore) AppendFrequency(runID model.ID, values []float64) {
	s.mu.Lock()
	trace := s.traces[runID]
	trace.RunID = runID
	trace.Frequencies = append(trace.Frequencies, values...)
	s.traces[runID] = trace
	s.mu.Unlock()
}

func (s *TraceStore) AppendThickness(sample model.ThicknessSample) {
	s.mu.Lock()
	trace := s.traces[sample.RunID]
	trace.RunID = sample.RunID
	trace.Thickness = append(trace.Thickness, sample)
	s.traces[sample.RunID] = trace
	s.mu.Unlock()
}

func (s *TraceStore) Trace(runID model.ID) Trace {
	s.mu.RLock()
	defer s.mu.RUnlock()
	source := s.traces[runID]
	result := Trace{RunID: source.RunID}
	result.Frequencies = append([]float64(nil), source.Frequencies...)
	result.Thickness = append([]model.ThicknessSample(nil), source.Thickness...)
	return result
}

func (s *TraceStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.traces)
}
