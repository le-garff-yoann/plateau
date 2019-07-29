package broadcaster

import (
	"sync"

	uuid "github.com/satori/go.uuid"
)

// Broadcaster delivers data from the emitter to his subscriber.
type Broadcaster struct {
	mux sync.RWMutex

	emitter     chan interface{}
	subscribers map[uuid.UUID]chan interface{}

	done chan int
}

// New returns a new `Broascaster`.
func New() *Broadcaster {
	return &Broadcaster{
		emitter:     make(chan interface{}),
		subscribers: make(map[uuid.UUID]chan interface{}),
		done:        make(chan int),
	}
}

// Submit send data to the registered subscribers.
func (s *Broadcaster) Submit(e interface{}) {
	s.emitter <- e
}

// Subscribe subscribes a new client for receiving data emitted with `Submit`.
func (s *Broadcaster) Subscribe() (<-chan interface{}, uuid.UUID) {
	s.mux.Lock()
	defer s.mux.Unlock()

	uuid := uuid.NewV4()

	s.subscribers[uuid] = make(chan interface{})

	return s.subscribers[uuid], uuid
}

// Unsubscribe unsubscribes the client based on the `UUID` which was returned by `Subscribe`.
func (s *Broadcaster) Unsubscribe(uuid uuid.UUID) bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	_, ok := s.subscribers[uuid]
	if ok {
		delete(s.subscribers, uuid)
	}

	return ok
}

// Run starts the broadcaster.
func (s *Broadcaster) Run() {
	for {
		select {
		case ec := <-s.emitter:
			go func() {
				s.mux.RLock()
				defer s.mux.RUnlock()

				for _, rs := range s.subscribers {
					go func(rs chan interface{}) {
						rs <- ec
					}(rs)
				}
			}()
		case <-s.done:
			return
		}
	}
}

// Done stops the broadcaster.
func (s *Broadcaster) Done() {
	s.done <- 0
}
