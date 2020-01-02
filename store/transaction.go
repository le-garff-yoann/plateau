package store

import (
	"plateau/protocol"
	"time"
)

// Transaction represents a specialized transaction which exposes requests specific to this project.
//
// It has `Commit()` and `Abort()` which allows
// an [ACID](https://en.wikipedia.org/wiki/ACID)-compliant implementation.
//	- Depending on the implementation, there is no guarantee that you will be able to
//	continue requesting in the same transaction after an error has been returned by a request.
//	This principle must be taken into account with the **systematic** use of `Abort()` after one of these errors.
type Transaction interface {
	Commit() error
	Abort() error

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
