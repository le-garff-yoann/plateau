package rockpaperscissors

import (
	"plateau/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGame(t *testing.T) {
	t.Parallel()

	g := Game{}
	g.Init()

	require.Equal(t, "rock–paper–scissors", g.Name())

	require.Equal(t, "https://en.wikipedia.org/wiki/Rock-paper-scissors", g.Description())

	require.Equal(t, uint(2), g.MinPlayers())
	require.Equal(t, uint(2), g.MaxPlayers())
}

func TestIsMatchValid(t *testing.T) {
	t.Parallel()

	g := Game{}
	g.Init()

	require.NoError(t, g.IsMatchValid(&protocol.Match{NumberOfPlayersRequired: 2}))

	require.Error(t, g.IsMatchValid(&protocol.Match{}))

	require.Error(t, g.IsMatchValid(&protocol.Match{NumberOfPlayersRequired: 0}))
	require.Error(t, g.IsMatchValid(&protocol.Match{NumberOfPlayersRequired: 1}))

	for i := 3; i < 10000; i++ {
		require.Error(t, g.IsMatchValid(&protocol.Match{NumberOfPlayersRequired: uint(i)}))
	}
}
