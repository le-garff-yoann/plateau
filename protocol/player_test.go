package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlayerString(t *testing.T) {
	player := Player{Name: "foo"}

	require.Equal(t, player.Name, player.String())
}
