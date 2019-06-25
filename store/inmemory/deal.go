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

// CreateDealsChangeIterator ...
func (s *Store) CreateDealsChangeIterator(id string) (store.DealChangeIterator, error) {
	itr := DealChangeIterator{dealChangesBroadcaster: s.dealChangesBroadcaster}

	itr.dealChangesBroadcasterChan, itr.dealChangesBroadcasterUUID = s.dealChangesBroadcaster.Subscribe()

	return &itr, nil
}

// DealChangeIterator implements `store.DealChangeIterator` interface.
type DealChangeIterator struct {
	dealChangesBroadcaster *broadcaster.Broadcaster

	dealChangesBroadcasterChan <-chan interface{}
	dealChangesBroadcasterUUID uuid.UUID
}

// Next implements `store.DealChangeIterator` interface.
func (s *DealChangeIterator) Next(dealChange *store.DealChange) bool {
	v, ok := <-s.dealChangesBroadcasterChan
	if !ok {
		return false
	}

	*dealChange = v.(store.DealChange)

	return true
}

// Close implements `store.DealChangeIterator` interface.
func (s *DealChangeIterator) Close() error {
	s.dealChangesBroadcaster.Unsubscribe(s.dealChangesBroadcasterUUID)

	return nil
}
