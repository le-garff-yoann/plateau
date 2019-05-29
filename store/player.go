package store

import "plateau/protocol"

// PlayerGameStore ...
type PlayerGameStore interface {
	List() (names []string, err error)
	Read(name string) (*protocol.Player, error)

	IncreaseWins(name string, increase uint) error
	IncreaseLoses(name string, increase uint) error
	IncreaseTies(name string, increase uint) error
}

// PlayerStore ...
type PlayerStore interface {
	PlayerGameStore

	Create(protocol.Player) error
}
