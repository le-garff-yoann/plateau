package postgresql

import "plateau/protocol"

// Deal ...
type Deal struct {
	ID uint

	Holder Player

	Messages []Message
}

func dealFromProtocolStruct(m *protocol.Deal) *Deal {
	var messages []Message
	for _, msg := range m.Messages {
		messages = append(messages, *messageFromProtocolStruct(&msg))
	}

	return &Deal{
		Holder:   *playerFromProtocolStruct(&m.Holder),
		Messages: messages,
	}
}

func (s *Deal) toProtocolStruct() *protocol.Deal {
	var messages []protocol.Message
	for _, msg := range s.Messages {
		messages = append(messages, *msg.toProtocolStruct())
	}

	return &protocol.Deal{
		Holder:   *s.Holder.toProtocolStruct(),
		Messages: messages,
	}
}
