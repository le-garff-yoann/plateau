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
		T:           t,
		Game:        &Game{},
		Match:       protocol.Match{NumberOfPlayersRequired: 2},
		PlayersName: []string{"foo", "bar"},
	}

	server.SetupTestMatchRuntime(t, testMatchRuntime)
	defer testMatchRuntime.Stop()

	testMatchRuntime.TestRequest("foo", protocol.ReqPlayerWantToJoin, protocol.ResOK)
	testMatchRuntime.TestRequest("bar", protocol.ReqPlayerWantToJoin, protocol.ResOK)
	testMatchRuntime.TestRequest("foo", protocol.ReqPlayerWantToStartTheMatch, protocol.ResOK)
	testMatchRuntime.TestRequest("foo", protocol.ReqPlayerAccepts, protocol.ResOK)
	testMatchRuntime.TestRequest("bar", protocol.ReqPlayerAccepts, protocol.ResOK)

	fight := func(reqA, reqB protocol.Request) *protocol.Match {
		// Force the creation of the first deal.
		testMatchRuntime.TestRequest("foo", protocol.ReqListRequests, protocol.ResOK)

		trn, err := testMatchRuntime.Store().BeginTransaction()
		require.NoError(t, err)

		match, err := trn.MatchRead(testMatchRuntime.Match.ID)
		require.NoError(t, err)
		trn.Abort()

		initialHolder := protocol.IndexDeals(match.Deals, 0).Holder
		require.NotNil(t, initialHolder)

		testMatchRuntime.TestRequest(initialHolder.Name, reqA, protocol.ResOK)

		trn, err = testMatchRuntime.Store().BeginTransaction()
		require.NoError(t, err)

		match, err = trn.MatchRead(testMatchRuntime.Match.ID)
		require.NoError(t, err)
		trn.Abort()

		currentDeal := protocol.IndexDeals(match.Deals, 0).WithMessagesConcealed(match.NextPlayer(initialHolder).Name)
		require.Empty(t, currentDeal.Messages[len(currentDeal.Messages)-1].Code)

		testMatchRuntime.TestRequest(match.NextPlayer(initialHolder).Name, reqB, protocol.ResOK)

		trn, err = testMatchRuntime.Store().BeginTransaction()
		require.NoError(t, err)

		match, err = trn.MatchRead(testMatchRuntime.Match.ID)
		require.NoError(t, err)
		trn.Abort()

		return match
	}

	require.False(t, fight(ReqRock, ReqRock).IsEnded())
	require.False(t, fight(ReqPaper, ReqPaper).IsEnded())
	require.False(t, fight(ReqScissors, ReqScissors).IsEnded())

	finalMatch := fight(ReqScissors, ReqRock)

	require.True(t, finalMatch.IsEnded())

	require.NotEqual(t, finalMatch.Players[0].Wins, finalMatch.Players[1].Wins)
	require.NotEqual(t, finalMatch.Players[0].Loses, finalMatch.Players[1].Loses)
}
