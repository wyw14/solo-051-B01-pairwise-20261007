package vacuum

import (
	"context"
	"github.com/wyw14/cry-161/internal/model"
	"sync"
	"time"
)

type Controller struct {
	mu       sync.RWMutex
	pressure float64
	pumping  bool
	stable   bool
}

func NewController() *Controller { return &Controller{pressure: 1.0} }

func (c *Controller) Start(ctx context.Context) error {
	c.mu.Lock()
	c.pumping = true
	c.stable = false
	c.mu.Unlock()
	select {
	case <-time.After(5 * time.Millisecond):
	case <-ctx.Done():
		return ctx.Err()
	}
	c.mu.Lock()
	c.pressure = 0.00001
	c.stable = true
	c.mu.Unlock()
	return nil
}
func (c *Controller) Stop()             { c.mu.Lock(); c.pumping = false; c.mu.Unlock() }
func (c *Controller) Pressure() float64 { c.mu.RLock(); defer c.mu.RUnlock(); return c.pressure }
func (c *Controller) Stable() bool      { c.mu.RLock(); defer c.mu.RUnlock(); return c.stable }
func (c *Controller) Snapshot(name string) model.ChamberSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return model.ChamberSnapshot{Name: name, Pressure: c.pressure, Stable: c.stable, UpdatedAt: time.Now().UTC()}
}
