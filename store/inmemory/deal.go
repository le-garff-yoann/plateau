package inmemory

import (
	"plateau/broadcaster"
	"plateau/protocol"
	"plateau/store"

	uuid "github.com/satori/go.uuid"
)

type deal struct {
	Holder protocol.Player

	Messages []protocol.Message
}

func dealFromProtocolStruct(m *protocol.Deal) *deal {
	return &deal{m.Holder, m.Messages}
}

func (s *deal) toProtocolStruct(pPlayers []*protocol.Player) *protocol.Deal {
	holder := s.Holder

	for _, p := range pPlayers {
		if p.Name == s.Holder.Name {
			holder = *p

			break
		}
	}

	return &protocol.Deal{
		Holder:   holder,
		Messages: s.Messages,
	}
}

// CreateDealsChangeIterator implements the `store.DealsChangeIterator` interface.
func (s *Store) CreateDealsChangeIterator(id string) (store.DealsChangeIterator, error) {
	itr := DealsChangeIterator{dealsChangeBroadcaster: s.dealsChangeBroadcaster}

	itr.dealsChangeBroadcasterChan, itr.dealsChangeBroadcasterUUID = s.dealsChangeBroadcaster.Subscribe()

	return &itr, nil
}

// DealsChangeIterator implements the `store.DealsChangeIterator` interface.
type DealsChangeIterator struct {
	dealsChangeBroadcaster *broadcaster.Broadcaster

	dealsChangeBroadcasterChan <-chan interface{}
	dealsChangeBroadcasterUUID uuid.UUID
}

// Next implements the `store.DealsChangeIterator` interface.
func (s *DealsChangeIterator) Next(dealChange *store.DealsChange) bool {
	v, ok := <-s.dealsChangeBroadcasterChan
	if !ok {
		return false
	}

	*dealChange = v.(store.DealsChange)

	return true
}

// Close implements the `store.DealsChangeIterator` interface.
func (s *DealsChangeIterator) Close() error {
	s.dealsChangeBroadcaster.Unsubscribe(s.dealsChangeBroadcasterUUID)

	return nil
}
