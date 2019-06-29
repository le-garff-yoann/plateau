package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchString(t *testing.T) {
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
