package rockpaperscissors

import (
	"plateau/protocol"
	"plateau/server"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGameRuntime(t *testing.T) {
	t.Parallel()

	testMatchRuntime := &server.TestMatchRuntime{
		Game:  &Game{},
		Match: protocol.Match{NumberOfPlayersRequired: 2},
		Players: []protocol.Player{
			protocol.Player{Name: "foo"},
			protocol.Player{Name: "bar"},
		},
	}

	server.SetupTestMatchRuntime(t, testMatchRuntime)

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

	fight := func(reqA, reqB protocol.Request) *protocol.Match {
		// Force the creation of the first deal.
		testMatchRuntime.ReqContainerHandlerFunc()(
			testMatchRuntime.Store().BeginTransaction(),
			&protocol.RequestContainer{
				Request: protocol.ReqListRequests,
				Player:  &testMatchRuntime.Players[0],
			})

		trn := testMatchRuntime.Store().BeginTransaction()

		match, _ := trn.MatchRead(testMatchRuntime.Match.ID)
		trn.Abort()

		initialHolder := protocol.IndexDeals(match.Deals, 0).Holder
		require.NotNil(t, initialHolder)

		require.Equal(t, protocol.ResOK, testMatchRuntime.ReqContainerHandlerFunc()(
			testMatchRuntime.Store().BeginTransaction(),
			&protocol.RequestContainer{
				Request: reqA,
				Player:  &initialHolder,
			}).Response,
		)

		require.Equal(t, protocol.ResOK, testMatchRuntime.ReqContainerHandlerFunc()(
			testMatchRuntime.Store().BeginTransaction(),
			&protocol.RequestContainer{
				Request: reqB,
				Player:  match.NextPLayer(initialHolder),
			}).Response,
		)

		trn = testMatchRuntime.Store().BeginTransaction()

		match, _ = trn.MatchRead(testMatchRuntime.Match.ID)
		trn.Abort()

		return match
	}

	require.Nil(t, fight(ReqRock, ReqRock).EndedAt)
	require.Nil(t, fight(ReqPaper, ReqPaper).EndedAt)
	require.Nil(t, fight(ReqScissors, ReqScissors).EndedAt)

	finalMatch := fight(ReqScissors, ReqRock)

	require.NotNil(t, finalMatch.EndedAt)

	require.NotEqual(t, finalMatch.Players[0].Wins, finalMatch.Players[1].Wins)
	require.NotEqual(t, finalMatch.Players[0].Loses, finalMatch.Players[1].Loses)
}
