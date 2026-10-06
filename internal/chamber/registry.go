package chamber

import (
	"fmt"
	"github.com/wyw14/cry-161/internal/model"
	"sort"
	"sync"
)

type Registry struct {
	mu       sync.RWMutex
	chambers map[string]*Chamber
}

func NewRegistry(values ...*Chamber) *Registry {
	registry := &Registry{chambers: map[string]*Chamber{}}
	for _, value := range values {
		if value != nil {
			registry.chambers[value.Snapshot().Name] = value
		}
	}
	return registry
}

func (r *Registry) Add(value *Chamber) error {
	if value == nil || value.Snapshot().Name == "" {
		return fmt.Errorf("chamber identity is required")
	}
	name := value.Snapshot().Name
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.chambers[name]; exists {
		return fmt.Errorf("chamber already exists")
	}
	r.chambers[name] = value
	return nil
}

func (r *Registry) Get(name string) (*Chamber, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, ok := r.chambers[name]
	return value, ok
}

func (r *Registry) Snapshots() []model.ChamberSnapshot {
	r.mu.RLock()
	values := make([]*Chamber, 0, len(r.chambers))
	for _, value := range r.chambers {
		values = append(values, value)
	}
	r.mu.RUnlock()
	result := make([]model.ChamberSnapshot, 0, len(values))
	for _, value := range values {
		result = append(result, value.Snapshot())
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}
