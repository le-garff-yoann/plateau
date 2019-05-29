package rockpaperscissors

import (
	"fmt"
	"plateau/protocol"
	"plateau/server"
)

// Game ...
type Game struct {
	name, description      string
	minPlayers, maxPlayers uint
}

// Init implements `server.Game` interface.
func (s *Game) Init() error {
	s.name = "rock–paper–scissors"

	s.minPlayers = 2
	s.maxPlayers = 2

	return nil
}

// Name implements `server.Game` interface.
func (s *Game) Name() string {
	return s.name
}

// Description implements `server.Game` interface.
func (s *Game) Description() string {
	return s.description
}

// IsMatchValid implements `server.Game` interface.
func (s *Game) IsMatchValid(g *protocol.Match) error {
	if !(g.NumberOfPlayersRequired >= s.MinPlayers() && g.NumberOfPlayersRequired <= s.MaxPlayers()) {
		return fmt.Errorf("The number of players must be between %d and %d", s.MinPlayers(), s.MaxPlayers())
	}

	return nil
}

// MinPlayers implements `server.Game` interface.
func (s *Game) MinPlayers() uint {
	return s.minPlayers
}

// MaxPlayers implements `server.Game` interface.
func (s *Game) MaxPlayers() uint {
	return s.maxPlayers
}

// Context implements `server.Game` interface.
func (s *Game) Context(matchRuntime *server.MatchRuntime, requestContainer *protocol.RequestContainer) *server.Context {
	return server.NewContext()
}
