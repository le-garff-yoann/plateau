package rockpaperscissors

import (
	"plateau/protocol"
)

// Game ...
type Game struct {
	name, description      string
	minPlayers, maxPlayers uint
}

// Init implements `server.Game` interface.
func (s *Game) Init() error {
	s.name = "rock–paper–scissors"

	s.description = "https://en.wikipedia.org/wiki/Rock-paper-scissors"

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
