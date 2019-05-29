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

		match  = protocol.Match{NumberOfPlayersRequired: 2}
		player = protocol.Player{Name: "foo"}
	)

	srv := &Server{
		game:          g,
		matchRuntimes: make(map[string]*MatchRuntime),
		store:         &inmemory.Store{},
	}

	srv.store.Open()

	id, _ := srv.store.Matchs().Create(match)

	srv.store.Players().Create(player)

	matchRuntime, err := srv.guardRuntime(id)
	require.NoError(t, err)

	require.Equal(t, protocol.ResOK, matchRuntime.requestContainerHandler(&protocol.RequestContainer{
		Request: protocol.ReqListRequests,
		Player:  &player,
	}).Response)

	require.Equal(t, protocol.ResOK, matchRuntime.requestContainerHandler(&protocol.RequestContainer{
		Request: protocol.ReqPlayerWantToJoin,
		Player:  &player,
	}).Response)

	require.Equal(t, protocol.ResForbidden, matchRuntime.requestContainerHandler(&protocol.RequestContainer{
		Request: protocol.ReqPlayerWantToStartTheGame,
		Player:  &player,
	}).Response)
}
