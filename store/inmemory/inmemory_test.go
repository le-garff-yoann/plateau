package inmemory

import (
	"plateau/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInMemory(t *testing.T) {
	t.Parallel()

	inm := &inMemory{}

	inm.players = append(inm.players, &protocol.Player{Name: "foo"})
	require.NotNil(t, inm.player("foo"))
	require.Nil(t, inm.player("bar"))

	inm.matchs = append(inm.matchs, &match{ID: "foo"})
	require.NotNil(t, inm.match("foo"))
	require.Nil(t, inm.match("bar"))
}
