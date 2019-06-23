package server

import (
	"plateau/protocol"
	"plateau/store/inmemory"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuntime(t *testing.T) {
	t.Parallel()

	var (
		g = &surrenderGame{}

		match   = protocol.Match{NumberOfPlayersRequired: 2}
		players = []protocol.Player{
			protocol.Player{Name: "foo"},
			protocol.Player{Name: "bar"},
		}
	)

	srv := &Server{
		game:          g,
		matchRuntimes: make(map[string]*MatchRuntime),
		store:         &inmemory.Store{},
	}

	srv.store.Open()

	id, _ := srv.store.Matchs().Create(match)

	for _, p := range players {
		srv.store.Players().Create(p)
	}

	matchRuntime, err := srv.guardRuntime(id)
	require.NoError(t, err)

	require.Equal(t, protocol.ResOK, matchRuntime.requestContainerHandler(&protocol.RequestContainer{
		Request: protocol.ReqListRequests,
		Player:  &players[0],
	}).Response)

	require.Equal(t, protocol.ResOK, matchRuntime.requestContainerHandler(&protocol.RequestContainer{
		Request: protocol.ReqPlayerWantToJoin,
		Player:  &players[0],
	}).Response)

	require.Equal(t, protocol.ResForbidden, matchRuntime.requestContainerHandler(&protocol.RequestContainer{
		Request: protocol.ReqPlayerWantToStartTheGame,
		Player:  &players[0],
	}).Response)

	require.Equal(t, protocol.ResOK, matchRuntime.requestContainerHandler(&protocol.RequestContainer{
		Request: protocol.ReqPlayerWantToJoin,
		Player:  &players[1],
	}).Response)

	require.Equal(t, protocol.ResOK, matchRuntime.requestContainerHandler(&protocol.RequestContainer{
		Request: protocol.ReqPlayerWantToStartTheGame,
		Player:  &players[0],
	}).Response)

	m, _ := srv.store.Matchs().Read(id)
	t.Log(m.Transactions)

	require.Equal(t, protocol.ResOK, matchRuntime.requestContainerHandler(&protocol.RequestContainer{
		Request: protocol.ReqPlayerAccepts,
		Player:  &players[0],
	}).Response)
	m, _ = srv.store.Matchs().Read(id)
	t.Log(m.Transactions)
	require.Equal(t, protocol.ResOK, matchRuntime.requestContainerHandler(&protocol.RequestContainer{
		Request: protocol.ReqPlayerAccepts,
		Player:  &players[1],
	}).Response)
}
