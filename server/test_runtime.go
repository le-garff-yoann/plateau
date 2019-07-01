package server

import (
	"plateau/protocol"
	"plateau/store"
	"plateau/store/inmemory"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestMatchRuntime ...
type TestMatchRuntime struct {
	T *testing.T

	Game

	srv *Server

	Match       protocol.Match
	PlayersName []string

	MatchRuntime *matchRuntime
}

// Store ...
func (s *TestMatchRuntime) Store() store.Store {
	return s.srv.store
}

// ReqContainerHandlerFunc ...
func (s *TestMatchRuntime) ReqContainerHandlerFunc() func(store.Transaction, *protocol.RequestContainer) *protocol.ResponseContainer {
	return s.MatchRuntime.reqContainerHandler
}

// TestRequest ...
func (s *TestMatchRuntime) TestRequest(playerName string, req protocol.Request, expectedRes protocol.Response) {
	require.Equal(s.T, expectedRes, s.ReqContainerHandlerFunc()(
		s.Store().BeginTransaction(),
		&protocol.RequestContainer{
			Request: req,
			Player:  &protocol.Player{Name: playerName},
		}).Response,
	)
}

// Close ...
func (s *TestMatchRuntime) Close() {
	s.srv.unguardRuntime(s.Match.ID)
	s.MatchRuntime = nil

	s.Store().Close()
}

// SetupTestMatchRuntime ...
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
