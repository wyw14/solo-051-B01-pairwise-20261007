package run

import (
	"github.com/wyw14/cry-161/internal/model"
	"sort"
	"sync"
	"time"
)

type TimelineEntry struct {
	ID     model.ID  `json:"id"`
	RunID  model.ID  `json:"run_id"`
	Name   string    `json:"name"`
	Detail string    `json:"detail"`
	At     time.Time `json:"at"`
}

type Timeline struct {
	mu      sync.RWMutex
	entries map[model.ID][]TimelineEntry
}

func NewTimeline() *Timeline {
	return &Timeline{entries: map[model.ID][]TimelineEntry{}}
}

func (t *Timeline) Add(runID model.ID, name string, detail string) TimelineEntry {
	entry := TimelineEntry{ID: model.NewID(), RunID: runID, Name: name, Detail: detail, At: time.Now().UTC()}
	t.mu.Lock()
	t.entries[runID] = append(t.entries[runID], entry)
	t.mu.Unlock()
	return entry
}

func (t *Timeline) Entries(runID model.ID) []TimelineEntry {
	t.mu.RLock()
	result := append([]TimelineEntry(nil), t.entries[runID]...)
	t.mu.RUnlock()
	sort.Slice(result, func(i, j int) bool { return result[i].At.Before(result[j].At) })
	return result
}

func (t *Timeline) Count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	total := 0
	for _, entries := range t.entries {
		total += len(entries)
	}
	return total
}
