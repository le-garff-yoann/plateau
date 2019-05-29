package store

import (
	"plateau/protocol"
	"time"
)

// MatchGameStore ...
type MatchGameStore interface {
	List() (ids []string, err error)
	Read(id string) (*protocol.Match, error)

	EndedAt(id string, val time.Time) error

	CreateTransaction(id string, transaction protocol.Transaction) error
	UpdateCurrentTransactionHolder(id, newHolderName string) error
	AddMessageToCurrentTransaction(id string, message protocol.Message) error
}

// MatchStore ...
type MatchStore interface {
	MatchGameStore

	Create(protocol.Match) (id string, err error)

	ConnectPlayer(id, playerName string) error
	DisconnectPlayer(id, playerName string) error
	PlayerJoins(id, playerName string) error
	PlayerLeaves(id, playerName string) error

	CreateTransactionsChangeIterator(id string) (TransactionChangeIterator, error)
}

// TransactionChange ...
type TransactionChange struct {
	Old *protocol.Transaction
	New *protocol.Transaction
}

// NewMessages ...
func (s *TransactionChange) NewMessages() (msgs []*protocol.Message) {
	for i := len(s.Old.Messages); i < len(s.New.Messages); i++ {
		msgs = append(msgs, &s.New.Messages[i])
	}

	return msgs
}

// TransactionChangeIterator ...
type TransactionChangeIterator interface {
	Next(*TransactionChange) bool
	Close() error
}
