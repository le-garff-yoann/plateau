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
	Game

	srv *Server

	Match   protocol.Match
	Players []protocol.Player

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

// Close ...
func (s *TestMatchRuntime) Close() {
	s.srv.unguardRuntime(s.Match.ID)
	s.MatchRuntime = nil

	s.Store().Close()
}

// SetupTestMatchRuntime ...
func SetupTestMatchRuntime(t *testing.T, testMatchRuntime *TestMatchRuntime) {
	testMatchRuntime.srv = &Server{
		game:          testMatchRuntime.Game,
		matchRuntimes: make(map[string]*matchRuntime),
		store:         &inmemory.Store{},
	}

	testMatchRuntime.Store().Open()

	trn := testMatchRuntime.Store().BeginTransaction()

	id, _ := trn.MatchCreate(testMatchRuntime.Match)

	for _, p := range testMatchRuntime.Players {
		trn.PlayerCreate(p)
	}

	m, _ := trn.MatchRead(id)
	testMatchRuntime.Match = *m

	trn.Commit()

	var err error
	testMatchRuntime.MatchRuntime, err = testMatchRuntime.srv.guardRuntime(id)
	require.NoError(t, err)

}
