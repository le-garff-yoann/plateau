package broadcaster

import (
	"sync"
)

// Broadcaster delivers data from the emitter to his subscriber.
type Broadcaster struct {
	mux sync.RWMutex

	emitter     chan interface{}
	subscribers map[chan<- interface{}]int

	done chan int
}

// New returns a new `Broascaster`.
func New() *Broadcaster {
	return &Broadcaster{
		emitter:     make(chan interface{}),
		subscribers: make(map[chan<- interface{}]int),
		done:        make(chan int),
	}
}

// Submit send data to the registered subscribers.
func (s *Broadcaster) Submit(e interface{}) {
	s.emitter <- e
}

// Register registers a new channel for receiving data emitted with `Submit`.
func (s *Broadcaster) Register(ch chan interface{}) bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	_, ok := s.subscribers[ch]
	if !ok {
		s.subscribers[ch] = 0
	}

	return !ok
}

// Unregister unregisters the given channel.
func (s *Broadcaster) Unregister(ch chan interface{}) bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	_, ok := s.subscribers[ch]
	if ok {
		delete(s.subscribers, ch)
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

				for rs := range s.subscribers {
					go func(rs chan<- interface{}) {
						rs <- ec
					}(rs)
				}
			}()
		case <-s.done:
			return
		}
	}
}

// Done stops the broadcaster and closes the emitter.
func (s *Broadcaster) Done() {
	s.done <- 0

	close(s.emitter)
}
