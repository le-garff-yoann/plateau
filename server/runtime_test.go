package server

import (
	"plateau/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGameRuntime(t *testing.T) {
	t.Parallel()

	testMatchRuntime := &TestMatchRuntime{
		Game:  &surrenderGame{},
		Match: protocol.Match{NumberOfPlayersRequired: 2},
		Players: []protocol.Player{
			protocol.Player{Name: "foo"},
			protocol.Player{Name: "bar"},
		},
	}

	SetupTestMatchRuntime(t, testMatchRuntime)

	require.Equal(t, protocol.ResOK, testMatchRuntime.ReqContainerHandlerFunc()(
		testMatchRuntime.Store().BeginTransaction(),
		&protocol.RequestContainer{
			Request: protocol.ReqListRequests,
			Player:  &testMatchRuntime.Players[0],
		}).Response,
	)

	require.Equal(t, protocol.ResOK, testMatchRuntime.ReqContainerHandlerFunc()(
		testMatchRuntime.Store().BeginTransaction(),
		&protocol.RequestContainer{
			Request: protocol.ReqPlayerWantToJoin,
			Player:  &testMatchRuntime.Players[0],
		}).Response,
	)

	require.Equal(t, protocol.ResForbidden, testMatchRuntime.ReqContainerHandlerFunc()(
		testMatchRuntime.Store().BeginTransaction(),
		&protocol.RequestContainer{
			Request: protocol.ReqPlayerWantToStartTheGame,
			Player:  &testMatchRuntime.Players[0],
		}).Response,
	)

	require.Equal(t, protocol.ResOK, testMatchRuntime.ReqContainerHandlerFunc()(
		testMatchRuntime.Store().BeginTransaction(),
		&protocol.RequestContainer{
			Request: protocol.ReqPlayerWantToJoin,
			Player:  &testMatchRuntime.Players[1],
		}).Response,
	)

	require.Equal(t, protocol.ResOK, testMatchRuntime.ReqContainerHandlerFunc()(
		testMatchRuntime.Store().BeginTransaction(),
		&protocol.RequestContainer{
			Request: protocol.ReqPlayerWantToStartTheGame,
			Player:  &testMatchRuntime.Players[0],
		}).Response,
	)

	require.Equal(t, protocol.ResOK, testMatchRuntime.ReqContainerHandlerFunc()(
		testMatchRuntime.Store().BeginTransaction(),
		&protocol.RequestContainer{
			Request: protocol.ReqPlayerAccepts,
			Player:  &testMatchRuntime.Players[0],
		}).Response,
	)

	require.Equal(t, protocol.ResOK, testMatchRuntime.ReqContainerHandlerFunc()(
		testMatchRuntime.Store().BeginTransaction(),
		&protocol.RequestContainer{
			Request: protocol.ReqPlayerAccepts,
			Player:  &testMatchRuntime.Players[1],
		}).Response,
	)
}
