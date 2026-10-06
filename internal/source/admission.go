package source

import (
	"context"
	"fmt"
	"sync"
)

type Admission struct {
	mu       sync.Mutex
	capacity int
	used     int
}

func NewAdmission(capacity int) *Admission { return &Admission{capacity: capacity} }
func (a *Admission) Acquire(ctx context.Context, weight int) error {
	if weight <= 0 {
		return fmt.Errorf("invalid power weight")
	}
	for {
		a.mu.Lock()
		if a.used+weight <= a.capacity {
			a.used += weight
			a.mu.Unlock()
			return nil
		}
		a.mu.Unlock()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}
func (a *Admission) Release(weight int) {
	a.mu.Lock()
	a.used -= weight
	if a.used < 0 {
		a.used = 0
	}
	a.mu.Unlock()
}
func (a *Admission) Used() int { a.mu.Lock(); defer a.mu.Unlock(); return a.used }
