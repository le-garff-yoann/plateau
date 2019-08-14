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

func TestNotificationsChannel(t *testing.T) {
	t.Parallel()

	s := &Store{}

	require.NoError(t, s.Open())
	defer func() {
		require.NoError(t, s.Close())
	}()

	ch := make(chan interface{})

	require.NoError(t, s.RegisterNotificationsChannel(ch))
	require.NoError(t, s.RegisterNotificationsChannel(ch))
}
