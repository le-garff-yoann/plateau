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
	require.Equal(s.T, expectedRes, s.reqContainerHandlerFunc()(
		s.Store().BeginTransaction(),
		&protocol.RequestContainer{
			Request: req,
			Player:  &protocol.Player{Name: playerName},
		}).Response,
	)
}

// Stop destroys the `matchRuntime` and stop the `Server`.
func (s *TestMatchRuntime) Stop() {
	s.srv.unguardRuntime(s.Match.ID)
	s.MatchRuntime = nil

	s.srv.Stop()
}

// SetupTestMatchRuntime initializes *testMatchRuntime* with a `Server`,
// a `protocol.Match`, a slice of `protocol.Player` and a `matchRuntime`.
func SetupTestMatchRuntime(t *testing.T, testMatchRuntime *TestMatchRuntime) {
	var err error
	testMatchRuntime.srv, err = New(testMatchRuntime.Game, &inmemory.Store{})
	require.NoError(t, err)

	trn := testMatchRuntime.Store().BeginTransaction()

	id, _ := trn.MatchCreate(testMatchRuntime.Match)

	for _, pName := range testMatchRuntime.PlayersName {
		trn.PlayerCreate(protocol.Player{Name: pName})
	}

	m, _ := trn.MatchRead(id)
	testMatchRuntime.Match = *m

	trn.Commit()

	testMatchRuntime.MatchRuntime, err = testMatchRuntime.srv.guardRuntime(id)
	require.NoError(t, err)
}
