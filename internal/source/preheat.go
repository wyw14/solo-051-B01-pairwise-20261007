package source

import (
	"context"
	"fmt"
	"time"
)

type Heater struct {
	StableAfter int
	Attempts    int
}

func (h *Heater) Preheat(ctx context.Context) error {
	h.Attempts++
	select {
	case <-time.After(time.Millisecond):
	case <-ctx.Done():
		return ctx.Err()
	}
	if h.StableAfter > 0 && h.Attempts <= h.StableAfter {
		return fmt.Errorf("source temperature unstable")
	}
	return nil
}
