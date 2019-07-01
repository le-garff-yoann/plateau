package server

import (
	"plateau/store/inmemory"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServerStartStop(t *testing.T) {
	t.Parallel()

	srv, err := New(&surrenderGame{}, &inmemory.Store{})
	require.NoError(t, err)

	require.NoError(t, srv.Start())
	require.NoError(t, srv.Stop())
}
