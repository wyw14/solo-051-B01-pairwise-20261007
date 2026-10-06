package source

import (
	"context"
	"time"
)

type Cooldown struct {
	controller *Controller
	interval   time.Duration
}

func NewCooldown(controller *Controller, interval time.Duration) *Cooldown {
	if interval <= 0 {
		interval = time.Millisecond
	}
	return &Cooldown{controller: controller, interval: interval}
}

func (c *Cooldown) Run(ctx context.Context) error {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		if c.controller.Power() <= 0 {
			return nil
		}
		c.controller.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
