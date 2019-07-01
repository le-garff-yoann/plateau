package server

import (
	"fmt"
	"plateau/protocol"
	"plateau/store"
)

type surrenderGame struct {
	name, description      string
	minPlayers, maxPlayers uint
}

// Init implements `Game` interface.
func (s *surrenderGame) Init() error {
	s.name = "surrender"

	s.minPlayers = 2
	s.maxPlayers = 2

	return nil
}

// Name implements `Game` interface.
func (s *surrenderGame) Name() string {
	return s.name
}

// Description implements `Game` interface.
func (s *surrenderGame) Description() string {
	return s.description
}

// IsMatchValid implements `Game` interface.
func (s *surrenderGame) IsMatchValid(g *protocol.Match) error {
	if !(g.NumberOfPlayersRequired >= s.MinPlayers() && g.NumberOfPlayersRequired <= s.MaxPlayers()) {
		return fmt.Errorf("The number of players must be between %d and %d", s.MinPlayers(), s.MaxPlayers())
	}

	return nil
}

// MinPlayers implements `Game` interface.
func (s *surrenderGame) MinPlayers() uint {
	return s.minPlayers
}

// MaxPlayers implements `Game` interface.
func (s *surrenderGame) MaxPlayers() uint {
	return s.maxPlayers
}

func (s *surrenderGame) Context(trn store.Transaction, reqContainer *protocol.RequestContainer) *Context {
	return NewContext()
}
