package run

import (
	"context"
	"github.com/wyw14/cry-161/internal/model"
	"github.com/wyw14/cry-161/internal/source"
	"sync"
)

type LayerCoordinator struct {
	Source    *source.Controller
	mu        sync.Mutex
	completed int
}

func (c *LayerCoordinator) Execute(ctx context.Context, layers []model.Layer) error {
	var wg sync.WaitGroup
	var errMu sync.Mutex
	var first error
	for _, layer := range layers {
		wg.Add(1)
		go func(item model.Layer) {
			defer wg.Done()
			if err := c.Source.Layer(ctx, item); err != nil {
				errMu.Lock()
				if first == nil {
					first = err
				}
				errMu.Unlock()
			}
			c.mu.Lock()
			c.completed++
			c.mu.Unlock()
		}(layer)
	}
	wg.Wait()
	return first
}
func (c *LayerCoordinator) Completed() int { c.mu.Lock(); defer c.mu.Unlock(); return c.completed }
