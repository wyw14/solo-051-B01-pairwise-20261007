package optical

import (
	"context"
	"fmt"
	"github.com/wyw14/cry-161/internal/model"
	"time"
)

type ReconnectPolicy struct {
	Attempts int
	Backoff  time.Duration
}

func DefaultReconnectPolicy() ReconnectPolicy {
	return ReconnectPolicy{Attempts: 3, Backoff: time.Millisecond}
}

func (p ReconnectPolicy) Connect(ctx context.Context, consumer *Consumer) (chan model.ThicknessSample, error) {
	if consumer == nil {
		return nil, fmt.Errorf("consumer is required")
	}
	if p.Attempts < 1 || p.Backoff < 0 {
		return nil, fmt.Errorf("reconnect policy is invalid")
	}
	for attempt := 0; attempt < p.Attempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		channel := consumer.Reconnect()
		if channel != nil {
			return channel, nil
		}
		timer := time.NewTimer(p.Backoff)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
	return nil, fmt.Errorf("optical reconnect exhausted")
}
