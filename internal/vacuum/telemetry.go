package vacuum

import (
	"sort"
	"sync"
	"time"
)

type Reading struct {
	Pressure float64   `json:"pressure"`
	Pumping  bool      `json:"pumping"`
	At       time.Time `json:"at"`
}

type Telemetry struct {
	mu       sync.RWMutex
	limit    int
	readings []Reading
}

func NewTelemetry(limit int) *Telemetry {
	if limit < 1 {
		limit = 1
	}
	return &Telemetry{limit: limit, readings: make([]Reading, 0, limit)}
}

func (t *Telemetry) Record(pressure float64, pumping bool) Reading {
	reading := Reading{Pressure: pressure, Pumping: pumping, At: time.Now().UTC()}
	t.mu.Lock()
	t.readings = append(t.readings, reading)
	if len(t.readings) > t.limit {
		t.readings = append([]Reading(nil), t.readings[len(t.readings)-t.limit:]...)
	}
	t.mu.Unlock()
	return reading
}

func (t *Telemetry) Readings() []Reading {
	t.mu.RLock()
	result := append([]Reading(nil), t.readings...)
	t.mu.RUnlock()
	sort.Slice(result, func(i, j int) bool { return result[i].At.Before(result[j].At) })
	return result
}

func (t *Telemetry) Stable(target float64, count int) bool {
	readings := t.Readings()
	if count < 1 || len(readings) < count {
		return false
	}
	for _, reading := range readings[len(readings)-count:] {
		if reading.Pressure > target || !reading.Pumping {
			return false
		}
	}
	return true
}
