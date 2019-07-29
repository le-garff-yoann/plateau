package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"plateau/protocol"
	"plateau/store/inmemory"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestGetMatchIDsHandlerHandler(t *testing.T) {
	t.Parallel()

	srv, err := Init("", "", &surrenderGame{}, &inmemory.Store{})
	require.NoError(t, err)

	require.NoError(t, srv.store.Open())
	defer func() {
		require.NoError(t, srv.store.Close())
	}()

	h := http.Handler(srv.router.Get("getMatchIDs").GetHandler())

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

	trn := srv.store.BeginTransaction()

	id, err := trn.MatchCreate(protocol.Match{})
	require.NoError(t, err)
	trn.Commit()

	rr = newRecorder()
	require.Equal(t, http.StatusOK, rr.Code)
	require.JSONEq(t, fmt.Sprintf(`["%s"]`, id), rr.Body.String())
}

func testCreateAndReadMatchHandler(t *testing.T, srv *Server) *protocol.Match {
	var (
		createH = http.Handler(srv.router.Get("createMatch").GetHandler())
		readH   = http.Handler(srv.router.Get("readMatch").GetHandler())
	)

	req, err := http.NewRequest("GET", "", nil)
	require.NoError(t, err)

	req = mux.SetURLVars(req, map[string]string{
		"id": "hopeitdoesnotexist",
	})

	newRecorder := func(h http.Handler, req *http.Request) *httptest.ResponseRecorder {
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		return rr
	}

	require.Equal(t, http.StatusNotFound, newRecorder(readH, req).Code)

	var match protocol.Match
	for i := 0; i < 4; i++ {
		req, err := http.NewRequest("POST", "", strings.NewReader(fmt.Sprintf(`{"number_of_players_required":%d}`, i)))
		require.NoError(t, err)

		res := newRecorder(createH, req)

		if i == 2 {
			require.Equal(t, http.StatusCreated, res.Code)

			require.NoError(t, json.NewDecoder(res.Body).Decode(&match))
		} else {
			require.Equal(t, http.StatusBadRequest, res.Code)
		}
	}

	req, err = http.NewRequest("GET", "", nil)
	require.NoError(t, err)

	req = mux.SetURLVars(req, map[string]string{
		"id": match.ID,
	})

	require.Equal(t, http.StatusOK, newRecorder(readH, req).Code)

	return &match
}

func TestCreateAndReadMatchHandler(t *testing.T) {
	t.Parallel()

	srv, err := Init("", "", &surrenderGame{}, &inmemory.Store{})
	require.NoError(t, err)

	require.NoError(t, srv.store.Open())
	defer func() {
		require.NoError(t, srv.store.Close())
	}()

	testCreateAndReadMatchHandler(t, srv)
}

func TestMatchPlayersNameHandler(t *testing.T) {
	t.Parallel()

	srv, err := Init("", "", &surrenderGame{}, &inmemory.Store{})
	require.NoError(t, err)

	require.NoError(t, srv.store.Open())
	defer func() {
		require.NoError(t, srv.store.Close())
	}()

	var (
		h = http.Handler(srv.router.Get("getMatchPlayersName").GetHandler())

		match  = testCreateAndReadMatchHandler(t, srv)
		player = protocol.Player{Name: "foo", Password: "foo"}
	)

	testRegisterAndLoginHandlers(t, srv, player.Name, player.Password)

	trn := srv.store.BeginTransaction()

	require.NoError(t, trn.MatchPlayerJoins(match.ID, player.Name))
	trn.Commit()

	req, err := http.NewRequest("GET", "", nil)
	require.NoError(t, err)

	req = mux.SetURLVars(req, map[string]string{
		"id": match.ID,
	})

	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.JSONEq(t, fmt.Sprintf(`["%s"]`, player.Name), rr.Body.String())
}

func TestMatchDealsHandler(t *testing.T) {
	t.Parallel()

	srv, err := Init("", "", &surrenderGame{}, &inmemory.Store{})
	require.NoError(t, err)

	require.NoError(t, srv.store.Open())
	defer func() {
		require.NoError(t, srv.store.Close())
	}()

	var (
		h = http.Handler(srv.router.Get("getMatchDeals").GetHandler())

		match  = testCreateAndReadMatchHandler(t, srv)
		player = protocol.Player{Name: "foo", Password: "foo"}

		trn = srv.store.BeginTransaction()
	)

	require.NoError(t, trn.MatchCreateDeal(match.ID, protocol.Deal{}))
	trn.Commit()

	_, loginRecorder := testRegisterAndLoginHandlers(t, srv, player.Name, player.Password)

	req, err := http.NewRequest("GET", "", nil)
	require.NoError(t, err)

	for _, c := range loginRecorder.Result().Cookies() {
		req.AddCookie(c)
	}

	req = mux.SetURLVars(req, map[string]string{
		"id": match.ID,
	})

	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestStreamMatchDealsChangeHandler(t *testing.T) {
	t.Parallel()

	srv, err := Init("", "", &surrenderGame{}, &inmemory.Store{})
	require.NoError(t, err)

	require.NoError(t, srv.store.Open())
	defer func() {
		require.NoError(t, srv.store.Close())
	}()

	match := testCreateAndReadMatchHandler(t, srv)
	_, loginRecorder := testRegisterAndLoginHandlers(t, srv, "foo", "bar")

	// srv.Start()
	// defer srv.Stop()

	h := http.Handler(srv.router.Get("streamMatchDealsChange").GetHandler())

	req, err := http.NewRequest("GET", "", nil)
	require.NoError(t, err)

	req = mux.SetURLVars(req, map[string]string{
		"id": match.ID,
	})

	newRecorder := func(async bool) *httptest.ResponseRecorder {
		for _, c := range loginRecorder.Result().Cookies() {
			req.AddCookie(c)
		}

		rr := httptest.NewRecorder()

		if async {
			go h.ServeHTTP(rr, req)
		} else {
			h.ServeHTTP(rr, req)
		}

		return rr
	}

	// FIXME:
	//	- rr.Body does not "refresh" on flusher.Flush().
	//	- Test if []protocol.Deals are "WithMessagesConcealed()".

	// var (
	// 	rr = newRecorder(true)

	// 	hits = 0

	// 	done = make(chan int)
	// )

	// go func() {
	// 	delim := []byte{':', ' '}

	// 	for {
	// 		b, err := rr.Body.ReadBytes('\n')
	// 		switch err {
	// 		case nil:
	// 			panic(err)
	// 		case io.EOF:
	// 			done <- 0

	// 			return
	// 		}

	// 		if len(b) < 2 {
	// 			continue
	// 		}

	// 		spl := bytes.Split(b, delim)

	// 		if len(spl) < 2 {
	// 			continue
	// 		}

	// 		hits++
	// 		if hits == 2 {
	// 			done <- 0
	// 		}
	// 	}
	// }()

	// trn := srv.store.BeginTransaction()

	// trn.MatchCreateDeal(match.ID, protocol.Deal{})
	// trn.MatchAddMessageToCurrentDeal(match.ID, protocol.Message{})

	// trn.Commit()

	// <-done

	// require.Equal(t, 2, hits)

	require.Equal(t, http.StatusOK, newRecorder(true).Code)

	trn := srv.store.BeginTransaction()

	require.NoError(t, trn.MatchEndedAt(match.ID, time.Now()))
	trn.Commit()

	require.Equal(t, http.StatusGone, newRecorder(false).Code)
}

func TestPatchMatchHandler(t *testing.T) {
	t.Parallel()

	srv, err := Init("", "", &surrenderGame{}, &inmemory.Store{})
	require.NoError(t, err)

	require.NoError(t, srv.store.Open())
	defer func() {
		require.NoError(t, srv.store.Close())
	}()

	match := testCreateAndReadMatchHandler(t, srv)

	var (
		patchMatchtH = http.Handler(srv.router.Get("patchMatch").GetHandler())

		players = map[string]protocol.Player{
			"foo": protocol.Player{Name: "foo", Password: "foo"},
			"bar": protocol.Player{Name: "bar", Password: "bar"},
		}
	)

	cookies := make(map[string][]*http.Cookie)

	for _, p := range players {
		_, rr := testRegisterAndLoginHandlers(t, srv, p.Name, p.Password)

		cookies[p.Name] = rr.Result().Cookies()
	}

	patchMatchtRecorder := func(p *protocol.Player, post string) *httptest.ResponseRecorder {
		req, err := http.NewRequest("PATCH", "", strings.NewReader(post))
		require.NoError(t, err)

		req = mux.SetURLVars(req, map[string]string{
			"id": match.ID,
		})

		for _, c := range cookies[p.Name] {
			req.AddCookie(c)
		}

		rr := httptest.NewRecorder()

		patchMatchtH.ServeHTTP(rr, req)

		return rr
	}

	for _, p := range players {
		require.Equal(t, http.StatusBadRequest, patchMatchtRecorder(&p, "").Code)
		require.Equal(t, http.StatusOK, patchMatchtRecorder(&p, "{}").Code)
	}
}
