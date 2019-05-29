package rethinkdb

import (
	"plateau/protocol"
)

type transaction struct {
	Holder player `rethinkdb:"holder_id,reference" rethinkdb_ref:"id"`

	Messages []message `rethinkdb:"messages"`
}

func transactionFromProtocolStruct(trx *protocol.Transaction) *transaction {
	var messages []message

	for _, msg := range trx.Messages {
		messages = append(messages, *messageFromProtocolStruct(&msg))
	}

	return &transaction{
		Holder:   *playerFromProtocolStruct(&trx.Holder),
		Messages: messages,
	}
}

func (s *transaction) toProtocolStruct() *protocol.Transaction {
	var messages []protocol.Message

	for _, msg := range s.Messages {
		messages = append(messages, *msg.toProtocolStruct())
	}

	return &protocol.Transaction{
		Holder:   *s.Holder.toProtocolStruct(),
		Messages: messages,
	}
}
