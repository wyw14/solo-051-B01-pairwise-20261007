package quartz

import (
	"github.com/wyw14/cry-161/internal/model"
	"sync"
	"time"
)

type Receiver struct {
	mu        sync.Mutex
	buffer    []float64
	consumers []chan []float64
	stream    chan []float64
}

func NewReceiver() *Receiver {
	return &Receiver{buffer: make([]float64, 0, 64), stream: make(chan []float64, 8)}
}
func (r *Receiver) Subscribe() <-chan []float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch := make(chan []float64, 8)
	r.consumers = append(r.consumers, ch)
	return ch
}
func (r *Receiver) Receive(values []float64) {
	r.mu.Lock()
	r.buffer = append(r.buffer[:0], values...)
	block := append([]float64(nil), r.buffer...)
	for _, ch := range r.consumers {
		ch <- block
	}
	r.mu.Unlock()
}
func (r *Receiver) Stream() chan<- []float64 { return r.stream }
func (r *Receiver) Run(stop <-chan struct{}, emit func(model.ThicknessSample)) {
	for {
		select {
		case values, ok := <-r.stream:
			if !ok {
				return
			}
			r.Receive(values)
			if len(values) > 0 {
				emit(model.ThicknessSample{Frequency: values[0], Valid: true, At: time.Now().UTC()})
			}
		case <-stop:
			return
		}
	}
}
func (r *Receiver) CloseStream() { close(r.stream) }
