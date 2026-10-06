package shutter

import (
	"sync"
	"time"
)

type Feedback struct {
	mu      sync.RWMutex
	state   string
	changed time.Time
}

func NewFeedback() *Feedback { return &Feedback{state: "closed", changed: time.Now().UTC()} }
func (f *Feedback) Update(state string) {
	f.mu.Lock()
	f.state = state
	f.changed = time.Now().UTC()
	f.mu.Unlock()
}
func (f *Feedback) State() (string, time.Time) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.state, f.changed
}
