package capoeira

import "sync"

// ChannelTransport implements Transport for two parties using Go channels.
type ChannelTransport struct {
	locations []string
	channels  map[string]chan interface{}
	lock      sync.RWMutex
}

func NewChannelTransport(partyA, partyB string) *ChannelTransport {
	return &ChannelTransport{
		locations: []string{partyA, partyB},
		channels: map[string]chan interface{}{
			partyA: make(chan interface{}, 1),
			partyB: make(chan interface{}, 1),
		},
	}
}

func (t *ChannelTransport) Send(from, to string, data interface{}) {
	t.lock.RLock()
	ch, ok := t.channels[to]
	t.lock.RUnlock()
	if ok {
		ch <- data
	}
}

func (t *ChannelTransport) Receive(from, at string) interface{} {
	t.lock.RLock()
	ch, ok := t.channels[at]
	t.lock.RUnlock()
	if ok {
		return <-ch
	}
	return nil
}

func (t *ChannelTransport) Locations() []string {
	return t.locations
}
