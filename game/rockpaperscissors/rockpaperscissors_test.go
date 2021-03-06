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

	require.Nil(t, g.IsMatchValid(&protocol.Match{}))

	require.Equal(t, uint(2), g.MinPlayers())
	require.Equal(t, uint(2), g.MaxPlayers())
}
