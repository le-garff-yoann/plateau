package game

import "plateau/store"

// Game ...
type Game interface {
	IsMatchValid(*store.Match) error

	Init() error

	Name() (name string)

	MinPlayers() (minPlayers uint)
	MaxPlayers() (maxPlayers uint)

	OnEvent(*store.Match, *store.EventContainer) error
}
