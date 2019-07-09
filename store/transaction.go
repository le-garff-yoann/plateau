package store

import (
	"plateau/protocol"
	"time"
)

// Transaction represents a specialized transaction.
//
// It has `Commit()` and `Abort()` which allows
// an [ACID](https://en.wikipedia.org/wiki/ACID)-compliant implementation.
type Transaction interface {
	Commit()
	Abort()

	Closed() bool
	Errors() []error

	PlayerList() (names []string, err error)
	PlayerCreate(protocol.Player) error
	PlayerRead(name string) (*protocol.Player, error)
	PlayerIncreaseWins(name string, increase uint) error
	PlayerIncreaseLoses(name string, increase uint) error
	PlayerIncreaseTies(name string, increase uint) error

	MatchList() (ids []string, err error)
	MatchCreate(protocol.Match) (id string, err error)
	MatchRead(id string) (*protocol.Match, error)
	MatchEndedAt(id string, val time.Time) error
	MatchPlayerJoins(id, playerName string) error
	MatchPlayerLeaves(id, playerName string) error
	MatchCreateDeal(id string, deal protocol.Deal) error
	MatchUpdateCurrentDealHolder(id, newHolderName string) error
	MatchAddMessageToCurrentDeal(id string, message protocol.Message) error
}

// TransactionScope tells the store the scope of
// the `Transaction` to optimize its isolation.
type TransactionScope struct {
	Mode TransactionScopeMode

	Subject interface{}
}

// IsSubjectAll returns `true` if the *Subject*
// is the entire store. Otherwise `false`.
func (s *TransactionScope) IsSubjectAll() bool {
	return s.Subject == nil
}

// IsSubjectPlayer returns non-nil if
// the *Subject* is a `protocol.Player`.
func (s *TransactionScope) IsSubjectPlayer() *protocol.Player {
	v, ok := s.Subject.(protocol.Player)
	if !ok {
		return nil
	}

	return &v
}

// IsSubjectMatch returns non-nil if
// the *Subject* is a `protocol.Match`.
func (s *TransactionScope) IsSubjectMatch() *protocol.Match {
	v, ok := s.Subject.(protocol.Match)
	if !ok {
		return nil
	}

	return &v
}

// TransactionScopeMode indicates whether the
// `Transaction` is read-only or read/write.
type TransactionScopeMode int

const (
	// TSReadMode indicates that the
	// `Transaction` will be read-only.
	TSReadMode = iota

	// TSReadWriteMode indicates that the
	// `Transaction` will be read/write.
	TSReadWriteMode
)
