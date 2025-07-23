package inmemory

import (
	"testing"

	"github.com/le-garff-yoann/plateau/protocol"
	"github.com/stretchr/testify/require"
)

func TestInMemory(t *testing.T) {
	t.Parallel()

	inm := &inMemory{}

	inm.Players = append(inm.Players, &protocol.Player{Name: "foo"})
	require.NotNil(t, inm.Player("foo"))
	require.Nil(t, inm.Player("bar"))

	inm.Matchs = append(inm.Matchs, &match{ID: "foo"})
	require.NotNil(t, inm.Match("foo"))
	require.Nil(t, inm.Match("bar"))

	require.IsType(t, inm, inm.Copy())
}
