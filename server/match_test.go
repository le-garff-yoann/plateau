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
	require.Equal(t, 200, rr.Code)
	require.JSONEq(t, "null", rr.Body.String())

	trn := srv.store.BeginTransaction()

	id, _ := trn.MatchCreate(protocol.Match{})
	trn.Commit()

	rr = newRecorder()
	require.Equal(t, 200, rr.Code)
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

	require.Equal(t, 404, newRecorder(readH, req).Code)

	req, err = http.NewRequest("POST", "", strings.NewReader(`{"number_of_players_required":2}`))
	require.NoError(t, err)

	res := newRecorder(createH, req)
	require.Equal(t, 201, res.Code)

	var match protocol.Match
	require.NoError(t, json.NewDecoder(res.Body).Decode(&match))

	req, err = http.NewRequest("GET", "", nil)
	require.NoError(t, err)

	req = mux.SetURLVars(req, map[string]string{
		"id": match.ID,
	})

	require.Equal(t, 200, newRecorder(readH, req).Code)

	return srv, &match
}

func TestCreateAndReadMatchHandler(t *testing.T) {
	t.Parallel()

	testCreateAndReadMatchHandler(t)
}

// func TestConnectMatchHandler(t *testing.T) {
// 	t.Parallel()

// 	srv, match := testCreateAndReadMatchHandler(t)

// 	srv.Start()
// 	defer srv.Stop()

// 	var (
// 		registerH = http.Handler(srv.router.Get("registerUser").GetHandler())
// 		loginH    = http.Handler(srv.router.Get("loginUser").GetHandler())

// 		connectH = http.Handler(srv.router.Get("connectMatch").GetHandler())

// 		players = map[string]protocol.Player{
// 			"foo": protocol.Player{Name: "foo", Password: "foo"},
// 			"bar": protocol.Player{Name: "bar", Password: "bar"},
// 		}
// 	)

// 	headers := make(map[string]http.Header)

// 	newPlayerRecorder := func(h http.Handler, p *protocol.Player) *httptest.ResponseRecorder {
// 		req, err := http.NewRequest("POST", "",
// 			strings.NewReader(fmt.Sprintf(`{"username":"%s","password":"%s"}`, p.Name, p.Password)))
// 		require.NoError(t, err)

// 		rr := httptest.NewRecorder()

// 		h.ServeHTTP(rr, req)

// 		require.Equal(t, 201, rr.Code)

// 		return rr
// 	}

// 	for _, p := range players {
// 		newPlayerRecorder(registerH, &p)

// 		rr := newPlayerRecorder(loginH, &p)

// 		headers[p.Name] = http.Header{}

// 		for _, c := range rr.Result().Cookies() {
// 			headers[p.Name].Add(c.Name, c.Value)
// 		}

// 		headers[p.Name].Add("X-Interactive", "true")
// 	}

// 	d := wstest.NewDialer(connectH)

// 	c, res, err := d.Dial(fmt.Sprintf("ws://x/api/matchs/%s", match.ID), headers["foo"])
// 	require.NoError(t, err)

// 	c.Close()

// 	require.Equal(t, http.StatusSwitchingProtocols, res.StatusCode)
// }
