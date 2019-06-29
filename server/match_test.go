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

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestGetMatchIDsHandlerHandler(t *testing.T) {
	t.Parallel()

	srv, err := New("", "", &surrenderGame{}, &inmemory.Store{})
	require.NoError(t, err)

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
	require.JSONEq(t, "null", rr.Body.String())

	trn := srv.store.BeginTransaction()

	id, _ := trn.MatchCreate(protocol.Match{})
	trn.Commit()

	rr = newRecorder()
	require.Equal(t, http.StatusOK, rr.Code)
	require.JSONEq(t, fmt.Sprintf(`["%s"]`, id), rr.Body.String())
}

func testCreateAndReadMatchHandler(t *testing.T) (*Server, *protocol.Match) {
	srv, err := New("", "", &surrenderGame{}, &inmemory.Store{})
	require.NoError(t, err)

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

	req, err = http.NewRequest("POST", "", strings.NewReader(`{"number_of_players_required":2}`))
	require.NoError(t, err)

	res := newRecorder(createH, req)
	require.Equal(t, http.StatusCreated, res.Code)

	var match protocol.Match
	require.NoError(t, json.NewDecoder(res.Body).Decode(&match))

	req, err = http.NewRequest("GET", "", nil)
	require.NoError(t, err)

	req = mux.SetURLVars(req, map[string]string{
		"id": match.ID,
	})

	require.Equal(t, http.StatusOK, newRecorder(readH, req).Code)

	return srv, &match
}

func TestCreateAndReadMatchHandler(t *testing.T) {
	t.Parallel()

	testCreateAndReadMatchHandler(t)
}

func TestStreamMatchNotificationsHandler(t *testing.T) {
	t.Parallel()

	// srv, match := testCreateAndReadMatchHandler(t)
	srv, _ := testCreateAndReadMatchHandler(t)

	// srv.Start()
	// defer srv.Stop()

	h := http.Handler(srv.router.Get("streamMatchNotifications").GetHandler())

	req, err := http.NewRequest("GET", "", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	go h.ServeHTTP(rr, req)

	// FIXME: rr.Body does not "refresh" on flusher.Flush().

	// var (
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

	require.Equal(t, http.StatusOK, rr.Code)
	// require.Equal(t, 2, hits)
}

func TestPatchMatchHandler(t *testing.T) {
	t.Parallel()

	srv, match := testCreateAndReadMatchHandler(t)

	var (
		registerH = http.Handler(srv.router.Get("registerUser").GetHandler())
		loginH    = http.Handler(srv.router.Get("loginUser").GetHandler())

		patchMatchtH = http.Handler(srv.router.Get("patchMatch").GetHandler())

		players = map[string]protocol.Player{
			"foo": protocol.Player{Name: "foo", Password: "foo"},
			"bar": protocol.Player{Name: "bar", Password: "bar"},
		}
	)

	cookies := make(map[string][]*http.Cookie)

	newPlayerRecorder := func(h http.Handler, p *protocol.Player) *httptest.ResponseRecorder {
		req, err := http.NewRequest("POST", "",
			strings.NewReader(fmt.Sprintf(`{"username":"%s","password":"%s"}`, p.Name, p.Password)))
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		require.Equal(t, http.StatusCreated, rr.Code)

		return rr
	}

	for _, p := range players {
		newPlayerRecorder(registerH, &p)

		rr := newPlayerRecorder(loginH, &p)

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
