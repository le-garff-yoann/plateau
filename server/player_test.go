package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"plateau/protocol"
	"plateau/store/inmemory"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestGetPlayersNameHandler(t *testing.T) {
	t.Parallel()

	srv, err := Init("", "", &surrenderGame{}, &inmemory.Store{})
	require.NoError(t, err)

	require.NoError(t, srv.store.Open())
	defer func() {
		require.NoError(t, srv.store.Close())
	}()

	h := http.Handler(srv.router.Get("getPlayersName").GetHandler())

	req, err := http.NewRequest("GET", "", nil)
	require.NoError(t, err)

	newRecorder := func() *httptest.ResponseRecorder {
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		return rr
	}

	rr := newRecorder()
	require.Equal(t, http.StatusOK, rr.Code)
	require.JSONEq(t, "[]", rr.Body.String())

	player := protocol.Player{Name: "foo"}

	trn, err := srv.store.BeginTransaction()
	require.NoError(t, err)

	require.NoError(t, trn.PlayerCreate(player))
	require.NoError(t, trn.Commit())

	rr = newRecorder()
	require.Equal(t, http.StatusOK, rr.Code)
	require.JSONEq(t, fmt.Sprintf(`["%s"]`, player.Name), rr.Body.String())
}

func TestReadPlayerNameHandler(t *testing.T) {
	t.Parallel()

	srv, err := Init("", "", &surrenderGame{}, &inmemory.Store{})
	require.NoError(t, err)

	require.NoError(t, srv.store.Open())
	defer func() {
		require.NoError(t, srv.store.Close())
	}()

	var (
		h = http.Handler(srv.router.Get("readPlayer").GetHandler())

		player = protocol.Player{Name: "foo"}
	)

	req, err := http.NewRequest("GET", "", nil)
	require.NoError(t, err)

	req = mux.SetURLVars(req, map[string]string{
		"name": player.Name,
	})

	newRecorder := func() *httptest.ResponseRecorder {
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		return rr
	}

	require.Equal(t, http.StatusNotFound, newRecorder().Code)

	trn, err := srv.store.BeginTransaction()
	require.NoError(t, err)

	require.NoError(t, trn.PlayerCreate(player))
	require.NoError(t, trn.Commit())

	require.Equal(t, http.StatusOK, newRecorder().Code)
}
