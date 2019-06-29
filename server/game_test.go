package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"plateau/protocol"
	"plateau/store"
	"plateau/store/inmemory"
	"testing"

	"github.com/stretchr/testify/require"
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

func (s *surrenderGame) Context(matchRuntime *MatchRuntime, trn store.Transaction, reqContainer *protocol.RequestContainer) *Context {
	trn.Commit()

	return NewContext()
}

func TestGetGameDefinitionHandler(t *testing.T) {
	t.Parallel()

	srv, err := New("", "", &surrenderGame{}, &inmemory.Store{})
	require.NoError(t, err)

	h := http.Handler(srv.router.Get("readGame").GetHandler())

	req, err := http.NewRequest("GET", "", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.JSONEq(t, fmt.Sprintf(
		`{"name":"%s","description":"%s","min_players":%d,"max_players":%d}`,
		srv.game.Name(), srv.game.Description(),
		srv.game.MinPlayers(), srv.game.MaxPlayers(),
	), rr.Body.String())
}
