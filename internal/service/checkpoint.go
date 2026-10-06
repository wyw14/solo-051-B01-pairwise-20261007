package service

import (
	"context"
	"fmt"
	"github.com/wyw14/cry-161/internal/journal"
	"github.com/wyw14/cry-161/internal/model"
	"os"
	"sort"
	"sync"
	"time"
)

type Checkpoint struct {
	ID        model.ID              `json:"id"`
	RunID     model.ID              `json:"run_id"`
	Sequence  int64                 `json:"sequence"`
	State     model.RunState        `json:"state"`
	Layer     int                   `json:"layer"`
	Chamber   model.ChamberSnapshot `json:"chamber"`
	Thickness model.ThicknessSample `json:"thickness"`
	CreatedAt time.Time             `json:"created_at"`
}

type CheckpointBook struct {
	mu        sync.RWMutex
	values    map[model.ID][]Checkpoint
	snapshots *journal.SnapshotStore
}

func NewCheckpointBook(dir string) *CheckpointBook {
	if dir == "" {
		dir = os.TempDir()
	}
	return &CheckpointBook{values: map[model.ID][]Checkpoint{}, snapshots: journal.NewSnapshotStore(dir)}
}

func (b *CheckpointBook) Capture(service *Service, runtime *Runtime, runID model.ID) (Checkpoint, error) {
	value, ok := service.Runs.Get(runID)
	if !ok {
		return Checkpoint{}, fmt.Errorf("run not found")
	}
	checkpoint := Checkpoint{ID: model.NewID(), RunID: runID, Sequence: runtime.NextSequence(), State: value.State, Layer: value.Layer, Chamber: service.Chamber.Snapshot(), Thickness: runtime.Monitor.Latest(), CreatedAt: time.Now().UTC()}
	b.mu.Lock()
	b.values[runID] = append(b.values[runID], checkpoint)
	b.mu.Unlock()
	snapshot := journal.Snapshot{Runs: service.Runs.List(), Recipes: service.Recipes.List(), Chamber: service.Chamber.Snapshot(), Incidents: runtime.Incidents.List(), Sequence: checkpoint.Sequence}
	if err := b.snapshots.Save(snapshot); err != nil {
		return checkpoint, err
	}
	return checkpoint, nil
}

func (b *CheckpointBook) List(runID model.ID) []Checkpoint {
	b.mu.RLock()
	result := append([]Checkpoint(nil), b.values[runID]...)
	b.mu.RUnlock()
	sort.Slice(result, func(i, j int) bool { return result[i].Sequence < result[j].Sequence })
	return result
}

func (b *CheckpointBook) Latest(runID model.ID) (Checkpoint, bool) {
	values := b.List(runID)
	if len(values) == 0 {
		return Checkpoint{}, false
	}
	return values[len(values)-1], true
}

func (b *CheckpointBook) Validate(checkpoint Checkpoint) error {
	if checkpoint.ID == "" || checkpoint.RunID == "" {
		return fmt.Errorf("checkpoint identity is required")
	}
	if checkpoint.Sequence < 1 || checkpoint.CreatedAt.IsZero() {
		return fmt.Errorf("checkpoint order is invalid")
	}
	if checkpoint.State == "" || checkpoint.Chamber.Name == "" {
		return fmt.Errorf("checkpoint state is incomplete")
	}
	return nil
}

func (b *CheckpointBook) Wait(ctx context.Context, runID model.ID, sequence int64) (Checkpoint, error) {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for {
		latest, ok := b.Latest(runID)
		if ok && latest.Sequence >= sequence {
			return latest, nil
		}
		select {
		case <-ctx.Done():
			return Checkpoint{}, ctx.Err()
		case <-ticker.C:
		}
	}
}

func (b *CheckpointBook) Restore() (journal.Snapshot, error) {
	snapshot, err := b.snapshots.Load()
	if err != nil {
		return journal.Snapshot{}, err
	}
	if snapshot.Sequence < 0 {
		return journal.Snapshot{}, fmt.Errorf("restored snapshot sequence is invalid")
	}
	return snapshot, nil
}
