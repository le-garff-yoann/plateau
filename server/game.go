package server

import (
	"net/http"
	"plateau/protocol"
	"plateau/server/response"
)

// Game ...
type Game interface {
	IsMatchValid(*protocol.Match) error

	Init() error

	Name() string
	Description() string

	MinPlayers() (minPlayers uint)
	MaxPlayers() (maxPlayers uint)

	Context(matchRuntime *MatchRuntime, requestContainer *protocol.RequestContainer) *Context
}

func (s *Server) getGameDefinitionHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"name":        s.game.Name(),
		"description": s.game.Description(),
		"min_players": s.game.MinPlayers(),
		"max_players": s.game.MaxPlayers(),
	})
}
