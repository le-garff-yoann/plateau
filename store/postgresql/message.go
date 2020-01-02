package postgresql

import (
	"plateau/protocol"
)

// Message ...
type Message struct {
	Code    string
	Payload interface{}

	AllowedNamesCode    []string
	AllowedNamesPayload []string
}

func messageFromProtocolStruct(m *protocol.Message) *Message {
	return &Message{
		Code:                m.Code.String(),
		Payload:             m.Payload,
		AllowedNamesCode:    m.AllowedNamesCode,
		AllowedNamesPayload: m.AllowedNamesPayload,
	}
}

func (s *Message) toProtocolStruct() *protocol.Message {
	msg := protocol.Message{
		Code:                protocol.MessageCode(s.Code),
		AllowedNamesCode:    s.AllowedNamesCode,
		AllowedNamesPayload: s.AllowedNamesPayload,
	}

	msg.EncodePayload(s.Payload)

	return &msg
}
