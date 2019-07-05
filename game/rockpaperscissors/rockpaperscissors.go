package rockpaperscissors

import (
	"plateau/protocol"
)

// Game implements the `server.Game` interface.
type Game struct {
	name, description      string
	minPlayers, maxPlayers uint
}

// Init implements the `server.Game` interface.
func (s *Game) Init() error {
	s.name = "rock–paper–scissors"

	s.description = "https://en.wikipedia.org/wiki/Rock-paper-scissors"

	s.minPlayers = 2
	s.maxPlayers = 2

	return nil
}

// Name implements the `server.Game` interface.
func (s *Game) Name() string {
	return s.name
}

// Description implements the `server.Game` interface.
func (s *Game) Description() string {
	return s.description
}

// IsMatchValid implements the`server.Game` interface.
func (s *Game) IsMatchValid(g *protocol.Match) error {
	return nil
}

// MinPlayers implements the `server.Game` interface.
func (s *Game) MinPlayers() uint {
	return s.minPlayers
}

// MaxPlayers implements the `server.Game` interface.
func (s *Game) MaxPlayers() uint {
	return s.maxPlayers
}
