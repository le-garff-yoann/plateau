package broadcaster

import (
	"sync"

	uuid "github.com/satori/go.uuid"
)

// Broadcaster ...
type Broadcaster struct {
	mux sync.Mutex

	emitter     chan interface{}
	subscribers map[uuid.UUID]chan interface{}

	done chan int
}

// New ...
func New() *Broadcaster {
	return &Broadcaster{
		emitter:     make(chan interface{}),
		subscribers: make(map[uuid.UUID]chan interface{}),
		done:        make(chan int),
	}
}

// Submit ...
func (s *Broadcaster) Submit(e interface{}) {
	s.emitter <- e
}

// Subscribe ...
func (s *Broadcaster) Subscribe() (<-chan interface{}, uuid.UUID) {
	s.mux.Lock()
	defer s.mux.Unlock()

	uuid := uuid.NewV4()

	s.subscribers[uuid] = make(chan interface{})

	return s.subscribers[uuid], uuid
}

// Unsubscribe ...
func (s *Broadcaster) Unsubscribe(uuid uuid.UUID) bool {
	_, ok := s.subscribers[uuid]

	if ok {
		s.mux.Lock()
		defer s.mux.Unlock()

		close(s.subscribers[uuid])
		delete(s.subscribers, uuid)
	}

	return ok
}

// Run ...
func (s *Broadcaster) Run() {
	for {
		select {
		case ec := <-s.emitter:
			var wg sync.WaitGroup

			s.mux.Lock()

			for _, rs := range s.subscribers {
				wg.Add(1)

				go func(rs chan interface{}) {
					rs <- ec

					wg.Done()
				}(rs)
			}

			wg.Wait()
			s.mux.Unlock()
		case <-s.done:
			return
		}
	}
}

// Done ...
func (s *Broadcaster) Done() {
	s.done <- 0
}
