package rethinkdb

import "plateau/protocol"

type message struct {
	MessageCode protocol.MessageCode `rethinkdb:"message_code"`
	Payload     interface{}          `rethinkdb:"payload"`
}

func messageFromProtocolStruct(msg *protocol.Message) *message {
	return &message{
		MessageCode: msg.MessageCode,
		Payload:     msg.Payload,
	}
}

func (s *message) toProtocolStruct() *protocol.Message {
	return &protocol.Message{
		MessageCode: s.MessageCode,
		Payload:     s.Payload,
	}
}
