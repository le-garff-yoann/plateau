package rockpaperscissors

import (
	"errors"
	"fmt"
	"plateau/event"
	"plateau/store"
)

// Game ...
type Game struct {
	name                   string
	minPlayers, maxPlayers uint
}

// Init ...
func (s *Game) Init() error {
	s.name = "rock–paper–scissors"

	s.minPlayers = 2
	s.maxPlayers = 2

	return nil
}

// Name implements `game.Game` interface.
func (s *Game) Name() string {
	return s.name
}

// IsMatchValid ...
func (s *Game) IsMatchValid(g *store.Match) error {
	if !(g.NumberOfPlayersRequired >= s.MinPlayers() && g.NumberOfPlayersRequired <= s.MaxPlayers()) {
		return fmt.Errorf("the number of players must be between %d and %d", s.MinPlayers(), s.MaxPlayers())
	}

	return nil
}

// MinPlayers implements `game.Game` interface.
func (s *Game) MinPlayers() uint {
	return s.minPlayers
}

// MaxPlayers implements `game.Game` interface.
func (s *Game) MaxPlayers() uint {
	return s.maxPlayers
}

// OnEvent implements `game.Game` interface.
func (s *Game) OnEvent(m *store.Match, ec *store.EventContainer) error {
	switch ec.Event {
	case event.EPlayerWantToJoin:
		if uint(len(m.Players)) >= s.MaxPlayers() {
			return errors.New("There are too many players in that game")
		}
	case event.EPlayerWantToSurrender:
		return errors.New("You cannont concede on this game")
	}

	return nil
}
