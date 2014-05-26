// Package hub provides a simple event dispatcher for publish/subscribe pattern.
package hub

import "sync"

type Kind int

// Event is an interface for published events.
type Event interface {
	Kind() Kind
}

// Hub is an event dispatcher, publishes events to the subscribers
// which are subscribed for a specific event type.
type Hub struct {
	subscribers map[Kind][]handler
	m           sync.RWMutex
	seq         uint64
}

// New returns pointer to a new Hub.
func New() *Hub {
	return &Hub{subscribers: make(map[Kind][]handler)}
}

// Subscribe registers f for the event of a specific kind.
func (h *Hub) Subscribe(kind Kind, f func(Event)) (cancel func()) {
	h.m.Lock()
	h.seq++
	id := h.seq
	h.subscribers[kind] = append(h.subscribers[kind], handler{id: id, f: f})
	h.m.Unlock()
	return func() {
		h.m.Lock()
		a := h.subscribers[kind]
		for i, h := range a {
			if h.id == id {
				a[i], a = a[len(a)-1], a[:len(a)-1]
				break
			}
		}
		if len(a) == 0 {
			delete(h.subscribers, kind)
		}
		h.m.Unlock()
	}
}

// Publish an event to the subscribers.
func (h *Hub) Publish(e Event) {
	h.m.RLock()
	if handlers, ok := h.subscribers[e.Kind()]; ok {
		for _, h := range handlers {
			h.f(e)
		}
	}
	h.m.RUnlock()
}

// DefaultHub is the default Hub used by Publish and Subscribe.
var DefaultHub = New()

// Subscribe registers f for the event of a specific kind in the DefaultHub.
func Subscribe(kind Kind, f func(Event)) (cancel func()) {
	return DefaultHub.Subscribe(kind, f)
}

// Publish an event to the subscribers in DefaultHub.
func Publish(e Event) {
	DefaultHub.Publish(e)
}

type handler struct {
	f  func(Event)
	id uint64
}
