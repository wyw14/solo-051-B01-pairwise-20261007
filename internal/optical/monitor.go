package optical

import (
	"github.com/wyw14/cry-161/internal/model"
	"sync"
	"time"
)

type Monitor struct {
	mu      sync.RWMutex
	samples []model.ThicknessSample
	latest  model.ThicknessSample
}

func NewMonitor() *Monitor { return &Monitor{samples: make([]model.ThicknessSample, 0, 64)} }
func (m *Monitor) Add(sample model.ThicknessSample) {
	m.mu.Lock()
	m.samples = append(m.samples, sample)
	m.latest = sample
	m.mu.Unlock()
}
func (m *Monitor) Latest() model.ThicknessSample { m.mu.RLock(); defer m.mu.RUnlock(); return m.latest }
func (m *Monitor) Samples() []model.ThicknessSample {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]model.ThicknessSample, len(m.samples))
	copy(out, m.samples)
	return out
}
func (m *Monitor) Estimate(freq, transmission float64) model.ThicknessSample {
	return model.ThicknessSample{Frequency: freq, Transmission: transmission, Nanometers: freq*0.42 + (1-transmission)*10, Valid: freq > 0 && transmission >= 0, At: time.Now().UTC()}
}
