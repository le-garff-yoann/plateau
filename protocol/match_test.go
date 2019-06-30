package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchString(t *testing.T) {
	t.Parallel()

	match := Match{ID: "foo"}

	require.Equal(t, match.ID, match.String())
}

func TestMatchIsFull(t *testing.T) {
	t.Parallel()

	match := Match{NumberOfPlayersRequired: 2}
	require.False(t, match.IsFull())

	match.Players = []Player{Player{}, Player{}}
	require.True(t, match.IsFull())
}

func TestNextPlayer(t *testing.T) {
	t.Parallel()

	match := Match{Players: []Player{
		Player{Name: "foo"},
		Player{Name: "bar"},
	}}

	require.Nil(t, match.NextPLayer(Player{Name: "baz"}))
	require.Equal(t, "bar", match.NextPLayer(Player{Name: "foo"}).Name)
	require.Equal(t, "foo", match.NextPLayer(Player{Name: "bar"}).Name)
}

func TestMatchRandomPlayer(t *testing.T) {
	t.Parallel()

	match := Match{Players: []Player{
		Player{Name: "foo"},
		Player{Name: "bar"},
	}}

	for i := 0; i < 10000; i++ {
		require.NotPanics(t, func() { match.RandomPlayer() })
	}
}
