package rotation

import (
	"github.com/wyw14/cry-161/internal/model"
	"sync"
	"time"
)

type Controller struct {
	mu       sync.Mutex
	rpm      float64
	run      model.ID
	spinning bool
}

func NewController() *Controller { return &Controller{} }
func (c *Controller) Start(run model.ID, rpm float64) {
	c.mu.Lock()
	c.run = run
	c.rpm = rpm
	c.spinning = true
	c.mu.Unlock()
}
func (c *Controller) Stop() { c.mu.Lock(); c.spinning = false; c.rpm = 0; c.mu.Unlock() }
func (c *Controller) Status() map[string]any {
	c.mu.Lock()
	defer c.mu.Unlock()
	return map[string]any{"run_id": c.run, "rpm": c.rpm, "spinning": c.spinning, "at": time.Now().UTC()}
}
