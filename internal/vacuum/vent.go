package vacuum

import "context"

func (c *Controller) Vent(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	c.mu.Lock()
	c.pressure = 1
	c.stable = false
	c.mu.Unlock()
	return nil
}
