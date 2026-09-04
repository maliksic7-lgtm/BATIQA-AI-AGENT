// Package events provides a tiny in-process pub/sub broker for SSE live updates.
package events

import "sync"

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type subscriber struct {
	id int64
	ch chan Event
}

// Broker fan-outs events to subscribers. Slow consumers drop messages
// rather than blocking publishers.
type Broker struct {
	mu   sync.Mutex
	next int64
	subs map[int64]subscriber
}

func New() *Broker {
	return &Broker{subs: map[int64]subscriber{}}
}

// Subscribe registers a new subscriber channel.
func (b *Broker) Subscribe() (int64, <-chan Event) {
	ch := make(chan Event, 16)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.next++
	id := b.next
	b.subs[id] = subscriber{id: id, ch: ch}
	return id, ch
}

// Unsubscribe removes a subscriber and closes its channel.
func (b *Broker) Unsubscribe(id int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if s, ok := b.subs[id]; ok {
		close(s.ch)
		delete(b.subs, id)
	}
}

// Publish sends an event to all subscribers; full channels skip the delivery.
func (b *Broker) Publish(eventType string, data interface{}) {
	ev := Event{Type: eventType, Data: data}
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, s := range b.subs {
		select {
		case s.ch <- ev:
		default: // slow consumer: drop instead of blocking
		}
	}
}
