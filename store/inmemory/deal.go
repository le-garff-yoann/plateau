package inmemory

import (
	"plateau/protocol"
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
