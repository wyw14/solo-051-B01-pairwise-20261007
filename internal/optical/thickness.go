package optical

import (
	"github.com/wyw14/cry-161/internal/model"
	"sync"
)

type ThicknessWindow struct {
	mu     sync.RWMutex
	values []model.ThicknessSample
	limit  int
}

func NewThicknessWindow(limit int) *ThicknessWindow {
	if limit < 1 {
		limit = 1
	}
	return &ThicknessWindow{limit: limit, values: make([]model.ThicknessSample, 0, limit)}
}
func (w *ThicknessWindow) Add(sample model.ThicknessSample) {
	w.mu.Lock()
	w.values = append(w.values, sample)
	if len(w.values) > w.limit {
		w.values = w.values[len(w.values)-w.limit:]
	}
	w.mu.Unlock()
}
func (w *ThicknessWindow) Values() []model.ThicknessSample {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]model.ThicknessSample, len(w.values))
	copy(out, w.values)
	return out
}
