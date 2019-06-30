package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"plateau/store/inmemory"
	"testing"

	"github.com/stretchr/testify/require"
)

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
