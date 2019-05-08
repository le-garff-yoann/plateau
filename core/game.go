package core

import (
	"plateau/model"
)

// Game ...
type Game interface {
	New() (*Game, error)

	Name(*Game) string

	IsValid(*Game, *model.Game) error

	MinPlayers(*Game) uint
	MaxPlayers(*Game) uint

	OnEvent(*Game, *model.EventContainer) error
}
