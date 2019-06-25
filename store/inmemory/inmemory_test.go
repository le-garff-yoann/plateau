package inmemory

import (
	"plateau/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInMemory(t *testing.T) {
	t.Parallel()

	inm := &inMemory{}

	inm.Players = append(inm.Players, &protocol.Player{Name: "foo"})
	require.NotNil(t, inm.player("foo"))
	require.Nil(t, inm.player("bar"))

	inm.Matchs = append(inm.Matchs, &match{ID: "foo"})
	require.NotNil(t, inm.match("foo"))
	require.Nil(t, inm.match("bar"))

	require.IsType(t, inm, inm.copy())
}
