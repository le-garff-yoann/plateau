package inmemory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	t.Parallel()

	store := &Store{}

	require.NoError(t, store.Open())
	require.NoError(t, store.Close())
}
