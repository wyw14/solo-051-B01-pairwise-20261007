package substrate

import (
	"fmt"
	"github.com/wyw14/cry-161/internal/model"
	"sync"
)

type Carrier struct {
	mu     sync.Mutex
	slots  map[string]model.ID
	loaded bool
}

func NewCarrier() *Carrier { return &Carrier{slots: map[string]model.ID{}} }
func (c *Carrier) Load(batch string, runID model.ID) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.loaded {
		return fmt.Errorf("carrier already loaded")
	}
	c.slots[batch] = runID
	c.loaded = true
	return nil
}
func (c *Carrier) Unload(batch string) {
	c.mu.Lock()
	delete(c.slots, batch)
	c.loaded = false
	c.mu.Unlock()
}
func (c *Carrier) Run(batch string) (model.ID, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	id, ok := c.slots[batch]
	return id, ok
}
func (c *Carrier) Loaded() bool { c.mu.Lock(); defer c.mu.Unlock(); return c.loaded }
