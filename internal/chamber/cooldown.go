package chamber

import (
	"context"
	"fmt"
	"time"
)

type CooldownPolicy struct {
	MaximumSourcePower float64
	MinimumDuration    time.Duration
	PollInterval       time.Duration
}

func DefaultCooldownPolicy() CooldownPolicy {
	return CooldownPolicy{MaximumSourcePower: 0.05, MinimumDuration: time.Millisecond, PollInterval: time.Millisecond}
}

func (p CooldownPolicy) Wait(ctx context.Context, chamber *Chamber) error {
	if chamber == nil {
		return fmt.Errorf("chamber is required")
	}
	if p.MinimumDuration < 0 || p.PollInterval <= 0 {
		return fmt.Errorf("cooldown policy is invalid")
	}
	started := time.Now()
	ticker := time.NewTicker(p.PollInterval)
	defer ticker.Stop()
	for {
		snapshot := chamber.Snapshot()
		elapsed := time.Since(started)
		if snapshot.SourcePower <= p.MaximumSourcePower && elapsed >= p.MinimumDuration {
			chamber.SetInterlock("cooling", true)
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
