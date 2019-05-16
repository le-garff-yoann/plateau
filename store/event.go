package store

import (
	"fmt"
	"plateau/event"
	"strings"
	"sync"

	uuid "github.com/satori/go.uuid"
)

// EventContainer ...
type EventContainer struct {
	event.Event `json:"event"`

	Emitter   *Player   `json:"emitter,omitempty"`
	Receivers []*Player `json:"receivers,omitempty"`
	Subjects  []*Player `json:"subjects,omitempty"`

	Payload map[string]interface{} `json:"payload"`
}

func (s *EventContainer) String() string {
	r := []string{fmt.Sprintf(`Event: "%s"`, s.Event)}

	if s.Emitter.Name != "" {
		r = append(r, fmt.Sprintf(`Emitter: "%s"`, s.Emitter.Name))
	}

	if len(s.Receivers) > 0 {
		var names []string
		for _, p := range s.Receivers {
			names = append(names, p.Name)
		}

		r = append(r, fmt.Sprintf(`Receivers: "%s"`, strings.Join(names, ", ")))
	}

	if len(s.Subjects) > 0 {
		var names []string
		for _, p := range s.Receivers {
			names = append(names, p.Name)
		}

		r = append(r, fmt.Sprintf(`Subjects: "%s"`, strings.Join(names, ", ")))
	}

	return strings.Join(r, " - ")
}

// IsLegal ...
func (s *EventContainer) IsLegal() bool {
	switch s.Event {
	case
		event.EIllegal, event.EInternalError,
		event.EGameWantToStart, event.EGameStarts, event.EGameEnds,
		event.EPlayerConnects, event.EPlayerDisconnects, event.EPlayerReconnects,
		event.EPlayerJoins, event.EPlayerLeaves, event.EPlayerSurrenders:
		return s.Emitter == nil
	case event.EListEvents, event.EPlayerWantToJoin, event.EPlayerWantToLeave, event.EPlayerWantToSurrender:
		return true
	}

	return false
}

// EventContainerBroadcaster ...
type EventContainerBroadcaster struct {
	mux sync.Mutex

	Emitter   chan EventContainer
	receivers map[uuid.UUID]chan EventContainer

	Done chan int
}

// NewEventContainerBroadcaster ...
func NewEventContainerBroadcaster() *EventContainerBroadcaster {
	return &EventContainerBroadcaster{
		Emitter:   make(chan EventContainer),
		receivers: make(map[uuid.UUID]chan EventContainer),
		Done:      make(chan int),
	}
}

// CountReceivers ...
func (s *EventContainerBroadcaster) CountReceivers() int {
	s.mux.Lock()
	defer s.mux.Unlock()

	return len(s.receivers)
}

// Recv ...
func (s *EventContainerBroadcaster) Recv() (<-chan EventContainer, uuid.UUID) {
	s.mux.Lock()
	defer s.mux.Unlock()

	uuid := uuid.NewV4()

	s.receivers[uuid] = make(chan EventContainer)

	return s.receivers[uuid], uuid
}

// RemoveRecv ...
func (s *EventContainerBroadcaster) RemoveRecv(uuid uuid.UUID) bool {
	_, ok := s.receivers[uuid]

	if ok {
		s.mux.Lock()
		defer s.mux.Unlock()

		close(s.receivers[uuid])
		delete(s.receivers, uuid)
	}

	return ok
}

// RemoveAllRecv ...
func (s *EventContainerBroadcaster) RemoveAllRecv() {
	s.mux.Lock()
	defer s.mux.Unlock()

	for _, rs := range s.receivers {
		close(rs)
	}

	s.receivers = make(map[uuid.UUID]chan EventContainer)
}

// Run ...
func (s *EventContainerBroadcaster) Run() {
	for {
		select {
		case ec := <-s.Emitter:
			var wg sync.WaitGroup

			s.mux.Lock()

			for _, rs := range s.receivers {
				wg.Add(1)

				go func(rs chan EventContainer) {
					rs <- ec

					wg.Done()
				}(rs)
			}

			wg.Wait()
			s.mux.Unlock()
		case <-s.Done:
			return
		}
	}
}
