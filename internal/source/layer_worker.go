package source

import (
	"context"
	"fmt"
	"github.com/wyw14/cry-161/internal/model"
	"sync"
	"time"
)

type LayerResult struct {
	Layer      model.Layer
	StartedAt  time.Time
	FinishedAt time.Time
	Error      string
}

type Worker struct {
	mu         sync.Mutex
	controller *Controller
	results    []LayerResult
}

func NewWorker(controller *Controller) *Worker {
	return &Worker{controller: controller, results: make([]LayerResult, 0)}
}

func (w *Worker) Execute(ctx context.Context, layer model.Layer) error {
	if layer.Number < 1 || layer.Material == "" {
		return fmt.Errorf("layer is invalid")
	}
	result := LayerResult{Layer: layer, StartedAt: time.Now().UTC()}
	err := w.controller.Layer(ctx, layer)
	result.FinishedAt = time.Now().UTC()
	if err != nil {
		result.Error = err.Error()
	}
	w.mu.Lock()
	w.results = append(w.results, result)
	w.mu.Unlock()
	return err
}

func (w *Worker) Results() []LayerResult {
	w.mu.Lock()
	defer w.mu.Unlock()
	result := make([]LayerResult, len(w.results))
	copy(result, w.results)
	return result
}
