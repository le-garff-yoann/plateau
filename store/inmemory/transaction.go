package inmemory

import (
	"plateau/protocol"
)

type transaction struct {
	Holder protocol.Player

	Messages []protocol.Message
}

func transactionFromProtocolStruct(m *protocol.Transaction) *transaction {
	return &transaction{m.Holder, m.Messages}
}

func (s *transaction) toProtocolStruct(pPlayers []*protocol.Player) *protocol.Transaction {
	holder := s.Holder

	for _, p := range pPlayers {
		if p.Name == s.Holder.Name {
			holder = *p

			break
		}
	}

	return &protocol.Transaction{
		Holder:   holder,
		Messages: s.Messages,
	}
}
