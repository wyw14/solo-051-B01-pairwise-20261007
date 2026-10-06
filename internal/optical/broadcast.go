package optical

import (
	"github.com/wyw14/cry-161/internal/model"
	"sync"
)

type Broadcaster struct {
	mu          sync.RWMutex
	subscribers map[int]chan model.ThicknessSample
	next        int
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{subscribers: map[int]chan model.ThicknessSample{}}
}
func (b *Broadcaster) Subscribe(size int) (int, <-chan model.ThicknessSample) {
	if size < 1 {
		size = 1
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	id := b.next
	b.next++
	ch := make(chan model.ThicknessSample, size)
	b.subscribers[id] = ch
	return id, ch
}
func (b *Broadcaster) Unsubscribe(id int) {
	b.mu.Lock()
	if ch, ok := b.subscribers[id]; ok {
		delete(b.subscribers, id)
		close(ch)
	}
	b.mu.Unlock()
}
func (b *Broadcaster) Publish(sample model.ThicknessSample) {
	b.mu.RLock()
	list := make([]chan model.ThicknessSample, 0, len(b.subscribers))
	for _, ch := range b.subscribers {
		list = append(list, ch)
	}
	b.mu.RUnlock()
	for _, ch := range list {
		select {
		case ch <- sample:
		default:
		}
	}
}
