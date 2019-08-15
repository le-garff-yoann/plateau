package server

import (
	"plateau/protocol"
	"plateau/store"
	"plateau/store/inmemory"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestMatchRuntime is used to test the runtime of a `protocol.Match`.
type TestMatchRuntime struct {
	srv *Server

	T *testing.T

	Game

	Match       protocol.Match
	PlayersName []string

	MatchRuntime *matchRuntime
}

func (s *TestMatchRuntime) reqContainerHandlerFunc() func(store.Transaction, *protocol.RequestContainer) *protocol.ResponseContainer {
	return s.MatchRuntime.reqContainerHandler
}

// Store exposes the `store.Store` interface of the instanciated `Server`.
func (s *TestMatchRuntime) Store() store.Store {
	return s.srv.store
}

// TestRequest sends *req* as player *playerName* and expects
// that the answer returned is equal to *expectedRes*.
func (s *TestMatchRuntime) TestRequest(playerName string, req protocol.Request, expectedRes protocol.Response) {
	trn, err := s.Store().BeginTransaction()
	require.NoError(s.T, err)

	require.Equal(s.T, expectedRes, s.reqContainerHandlerFunc()(
		trn,
		&protocol.RequestContainer{
			Request: req,
			Player:  &protocol.Player{Name: playerName},
		}).Response,
	)
}

// Stop destroys the `matchRuntime` and closes the `store.Store` of the `Server`.
func (s *TestMatchRuntime) Stop() {
	s.srv.unguardRuntime(s.Match.ID)
	s.MatchRuntime = nil

	require.NoError(s.T, s.Store().Close())
}

// SetupTestMatchRuntime initializes *testMatchRuntime* with a `Server`,
// an `inmemory.Store`, a `protocol.Match`, a slice of `protocol.Player`
// and a `matchRuntime`.
//
// Note that this function will not start the `Server`.
func SetupTestMatchRuntime(t *testing.T, testMatchRuntime *TestMatchRuntime) {
	var err error
	testMatchRuntime.srv, err = New(testMatchRuntime.Game, &inmemory.Store{})
	require.NoError(t, err)

	require.NoError(t, testMatchRuntime.Store().Open())

	trn, err := testMatchRuntime.Store().BeginTransaction()
	require.NoError(t, err)

	id, err := trn.MatchCreate(testMatchRuntime.Match)
	require.NoError(t, err)

	for _, pName := range testMatchRuntime.PlayersName {
		require.NoError(t, trn.PlayerCreate(protocol.Player{Name: pName}))
	}

	m, err := trn.MatchRead(id)
	require.NoError(t, err)

	testMatchRuntime.Match = *m

	trn.Commit()

	testMatchRuntime.MatchRuntime, err = testMatchRuntime.srv.guardRuntime(id)
	require.NoError(t, err)
}
