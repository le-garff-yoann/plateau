package server

import (
	"net/http"
	"plateau/protocol"
	"plateau/server/response"
	"plateau/store"
)

// Game represents the definition of a game.
//	- `IsMatchValid()` is called when creating the game through the API.
//	- `Init()` initializes the implementation of this interface.
//	This can be useful if it have values that need to be defined.
//	- `Context()` is called on every `protocol.RequestContainer` received if the
//	match is started and not finished.
type Game interface {
	IsMatchValid(*protocol.Match) error

	Init() error

	Name() string
	Description() string

	MinPlayers() (minPlayers uint)
	MaxPlayers() (maxPlayers uint)

	Context(store.Transaction, *protocol.RequestContainer) *Context
}

func (s *Server) getGameDefinitionHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"name":        s.game.Name(),
		"description": s.game.Description(),
		"min_players": s.game.MinPlayers(),
		"max_players": s.game.MaxPlayers(),
	})
}
