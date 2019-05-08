package game

import (
	"errors"
	"fmt"
	"plateau/model"
)

// Game ...
type Game struct {
	name       string
	minPlayers uint
	maxPlayers uint
}

// New ...
func New() (*Game, error) {
	return &Game{
		name:       "rock–paper–scissors",
		minPlayers: 2,
		maxPlayers: 2,
	}, nil
}

// Name implements `core.Game` interface.
func (s *Game) Name() string { return s.name }

// IsValid ...
func (s *Game) IsValid(g *model.Game) error {
	if !(g.NumberOfPlayersRequired >= s.MinPlayers() && g.NumberOfPlayersRequired <= s.MaxPlayers()) {
		return fmt.Errorf("The number of players must be between %d and %d", s.MinPlayers(), s.MaxPlayers())
	}

	return nil
}

// MinPlayers implements `core.Game` interface.
func (s *Game) MinPlayers() uint { return s.minPlayers }

// MaxPlayers implements `core.Game` interface.
func (s *Game) MaxPlayers() uint { return s.maxPlayers }

// OnEvents implements `core.Game` interface.
func (s *Game) OnEvent(c *model.EventContainer) error {
	switch c.Event {
	case model.EPlayerWantToJoin:
		if uint(len(c.Game.Players)) >= s.MaxPlayers() {
			return errors.New("There are too many players in that game")
		}
	case model.EPlayerWantToSurrender:
		return errors.New("You cannont concede on this game")
	}

	return nil
}
