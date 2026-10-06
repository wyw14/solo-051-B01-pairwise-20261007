package run

import (
	"context"
	"fmt"
	"github.com/wyw14/cry-161/internal/model"
	"sort"
	"sync"
	"time"
)

type ScheduledLayer struct {
	RunID      model.ID    `json:"run_id"`
	Layer      model.Layer `json:"layer"`
	Priority   int         `json:"priority"`
	EnqueuedAt time.Time   `json:"enqueued_at"`
	StartedAt  time.Time   `json:"started_at,omitempty"`
	FinishedAt time.Time   `json:"finished_at,omitempty"`
	Status     string      `json:"status"`
	Error      string      `json:"error,omitempty"`
}

type Scheduler struct {
	mu       sync.Mutex
	pending  []ScheduledLayer
	running  map[model.ID]ScheduledLayer
	finished []ScheduledLayer
	limit    int
}

func NewScheduler(limit int) *Scheduler {
	if limit < 1 {
		limit = 1
	}
	return &Scheduler{pending: make([]ScheduledLayer, 0), running: map[model.ID]ScheduledLayer{}, finished: make([]ScheduledLayer, 0), limit: limit}
}

func (s *Scheduler) Enqueue(runID model.ID, layer model.Layer, priority int) error {
	if runID == "" || layer.Number < 1 || layer.Material == "" {
		return fmt.Errorf("scheduled layer is invalid")
	}
	item := ScheduledLayer{RunID: runID, Layer: layer, Priority: priority, EnqueuedAt: time.Now().UTC(), Status: "pending"}
	s.mu.Lock()
	s.pending = append(s.pending, item)
	sort.SliceStable(s.pending, func(i, j int) bool {
		if s.pending[i].Priority == s.pending[j].Priority {
			return s.pending[i].EnqueuedAt.Before(s.pending[j].EnqueuedAt)
		}
		return s.pending[i].Priority > s.pending[j].Priority
	})
	s.mu.Unlock()
	return nil
}

func (s *Scheduler) Next() (ScheduledLayer, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.pending) == 0 || len(s.running) >= s.limit {
		return ScheduledLayer{}, false
	}
	item := s.pending[0]
	s.pending = append([]ScheduledLayer(nil), s.pending[1:]...)
	item.Status = "running"
	item.StartedAt = time.Now().UTC()
	s.running[item.RunID] = item
	return item, true
}

func (s *Scheduler) Finish(item ScheduledLayer, err error) ScheduledLayer {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.running, item.RunID)
	item.FinishedAt = time.Now().UTC()
	item.Status = "completed"
	if err != nil {
		item.Status = "failed"
		item.Error = err.Error()
	}
	s.finished = append(s.finished, item)
	return item
}

func (s *Scheduler) Execute(ctx context.Context, worker func(context.Context, model.ID, model.Layer) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		item, ok := s.Next()
		if !ok {
			pending, running, _ := s.Counts()
			if pending == 0 && running == 0 {
				return nil
			}
			time.Sleep(time.Millisecond)
			continue
		}
		err := worker(ctx, item.RunID, item.Layer)
		s.Finish(item, err)
		if err != nil {
			return err
		}
	}
}

func (s *Scheduler) Counts() (int, int, int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.pending), len(s.running), len(s.finished)
}

func (s *Scheduler) Results() []ScheduledLayer {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]ScheduledLayer, len(s.finished))
	copy(result, s.finished)
	return result
}
