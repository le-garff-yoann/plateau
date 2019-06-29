package store

import (
	"plateau/protocol"
	"time"
)

// Transaction ...
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
