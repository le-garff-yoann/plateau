package store

import (
	"plateau/model"
	"sync"
)

// ObjectStore ...
type ObjectStore interface {
	Mutex(*GameStore) sync.Mutex
}

// GameStore ...
type GameStore interface {
	ObjectStore

	Read(*GameStore) (model.Game, bool)
	Write(*GameStore, *model.Game) []error
}

// PlayerStore ...
type PlayerStore interface {
	ObjectStore

	Read(*PlayerStore) (model.Player, bool)
	Write(*PlayerStore, *model.Player) []error
}

// EventContainerStore ...
type EventContainerStore interface {
	ObjectStore

	Read(*EventContainerStore) (model.EventContainer, bool)
	Write(*EventContainerStore, *model.EventContainer) []error
}

// Store ...
type Store interface {
	Games() map[uint]*GameStore
	Players() map[string]*PlayerStore
	EventContainers() map[uint]*EventContainerStore
}

// Finish the postgres impl.
// The sessionStore should also have an interface.
// Test https://github.com/rethinkdb/rethinkdb-go
