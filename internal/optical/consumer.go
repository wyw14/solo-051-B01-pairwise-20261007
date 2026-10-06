package optical

import (
	"github.com/wyw14/cry-161/internal/model"
	"sync"
	"time"
)

type Consumer struct {
	mu         sync.Mutex
	samples    chan model.ThicknessSample
	monitor    *Monitor
	reconnects int
}

func NewConsumer(monitor *Monitor, samples chan model.ThicknessSample) *Consumer {
	return &Consumer{monitor: monitor, samples: samples}
}
func (c *Consumer) Run(stop <-chan struct{}) {
	for {
		c.mu.Lock()
		samples := c.samples
		c.mu.Unlock()
		if samples == nil {
			select {
			case <-stop:
				return
			case <-time.After(time.Millisecond):
				continue
			}
		}
		select {
		case sample, ok := <-samples:
			if !ok {
				c.mu.Lock()
				if c.samples == samples {
					c.samples = nil
				}
				c.mu.Unlock()
				continue
			}
			c.monitor.Add(sample)
		case <-stop:
			return
		}
	}
}
func (c *Consumer) Reconnect() chan model.ThicknessSample {
	c.mu.Lock()
	defer c.mu.Unlock()
	samples := make(chan model.ThicknessSample, 8)
	c.samples = samples
	c.reconnects++
	return samples
}
func (c *Consumer) ReconnectCount() int { c.mu.Lock(); defer c.mu.Unlock(); return c.reconnects }
func (c *Consumer) WaitForSample(timeout time.Duration) (model.ThicknessSample, bool) {
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	for {
		sample := c.monitor.Latest()
		if sample.At.IsZero() {
			select {
			case <-deadline.C:
				return model.ThicknessSample{}, false
			default:
				time.Sleep(time.Millisecond)
				continue
			}
		}
		return sample, true
	}
}
