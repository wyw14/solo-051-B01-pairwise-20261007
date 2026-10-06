package interlock

import "sync"

type Monitor struct {
	mu     sync.RWMutex
	alarms []string
}

func NewMonitor() *Monitor           { return &Monitor{alarms: make([]string, 0)} }
func (m *Monitor) Alarm(name string) { m.mu.Lock(); m.alarms = append(m.alarms, name); m.mu.Unlock() }
func (m *Monitor) Alarms() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]string, len(m.alarms))
	copy(out, m.alarms)
	return out
}
