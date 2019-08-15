package inmemory

import (
	"plateau/protocol"
	"plateau/store"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBeginTransactionCommit(t *testing.T) {
	t.Parallel()

	s := &Store{}

	require.NoError(t, s.Open())
	defer func() {
		require.NoError(t, s.Close())
	}()

	trn, err := s.BeginTransaction()
	require.NoError(t, err)

	id, err := trn.MatchCreate(protocol.Match{})
	require.NoError(t, err)

	_, err = trn.MatchRead(id)
	require.NoError(t, err)

	require.False(t, trn.Closed())
	require.NotPanics(t, func() { trn.Commit() })
	require.Panics(t, func() { trn.Commit() })
	require.Panics(t, func() { trn.Abort() })

	require.True(t, trn.Closed())

	trn, err = s.BeginTransaction()
	require.NoError(t, err)

	_, err = trn.MatchRead(id)
	require.NoError(t, err)
	require.NotPanics(t, func() { trn.Commit() })
}

func TestBeginTransactionAbort(t *testing.T) {
	t.Parallel()

	s := &Store{}

	require.NoError(t, s.Open())
	defer func() {
		require.NoError(t, s.Close())
	}()

	trn, err := s.BeginTransaction()
	require.NoError(t, err)

	id, err := trn.MatchCreate(protocol.Match{})
	require.NoError(t, err)

	_, err = trn.MatchRead(id)
	require.NoError(t, err)

	require.False(t, trn.Closed())
	trn.Abort()
	require.Panics(t, func() { trn.Commit() })
	require.Panics(t, func() { trn.Abort() })

	require.True(t, trn.Closed())

	require.Empty(t, trn.Errors())

	trn, err = s.BeginTransaction()
	require.NoError(t, err)

	_, err = trn.MatchRead(id)
	require.IsType(t, store.DontExistError(""), err)
}
