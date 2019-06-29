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

	trn := srv.store.BeginTransaction()

	id, _ := trn.MatchCreate(match)

	for _, p := range players {
		trn.PlayerCreate(p)
	}

	trn.Commit()

	matchRuntime, err := srv.guardRuntime(id)
	require.NoError(t, err)

	require.Equal(t, protocol.ResOK, matchRuntime.reqContainerHandler(
		srv.store.BeginTransaction(),
		&protocol.RequestContainer{
			Request: protocol.ReqListRequests,
			Player:  &players[0],
		}).Response,
	)

	require.Equal(t, protocol.ResOK, matchRuntime.reqContainerHandler(
		srv.store.BeginTransaction(),
		&protocol.RequestContainer{
			Request: protocol.ReqPlayerWantToJoin,
			Player:  &players[0],
		}).Response,
	)

	require.Equal(t, protocol.ResForbidden, matchRuntime.reqContainerHandler(
		srv.store.BeginTransaction(),
		&protocol.RequestContainer{
			Request: protocol.ReqPlayerWantToStartTheGame,
			Player:  &players[0],
		}).Response,
	)

	require.Equal(t, protocol.ResOK, matchRuntime.reqContainerHandler(
		srv.store.BeginTransaction(),
		&protocol.RequestContainer{
			Request: protocol.ReqPlayerWantToJoin,
			Player:  &players[1],
		}).Response,
	)

	require.Equal(t, protocol.ResOK, matchRuntime.reqContainerHandler(
		srv.store.BeginTransaction(),
		&protocol.RequestContainer{
			Request: protocol.ReqPlayerWantToStartTheGame,
			Player:  &players[0],
		}).Response,
	)

	require.Equal(t, protocol.ResOK, matchRuntime.reqContainerHandler(
		srv.store.BeginTransaction(),
		&protocol.RequestContainer{
			Request: protocol.ReqPlayerAccepts,
			Player:  &players[0],
		}).Response,
	)

	require.Equal(t, protocol.ResOK, matchRuntime.reqContainerHandler(
		srv.store.BeginTransaction(),
		&protocol.RequestContainer{
			Request: protocol.ReqPlayerAccepts,
			Player:  &players[1],
		}).Response,
	)
}
