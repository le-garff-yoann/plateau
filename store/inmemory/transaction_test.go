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
	s.Open()

	trn := s.BeginTransaction()

	id, err := trn.MatchCreate(protocol.Match{})
	require.NoError(t, err)

	_, err = trn.MatchRead(id)
	require.NoError(t, err)

	require.False(t, trn.Closed())
	trn.Commit()
	require.Panics(t, func() { trn.Commit() })
	require.Panics(t, func() { trn.Abort() })

	require.True(t, trn.Closed())

	trn = s.BeginTransaction()

	_, err = trn.MatchRead(id)
	require.NoError(t, err)
}

func TestBeginTransactionAbort(t *testing.T) {
	t.Parallel()

	s := &Store{}
	s.Open()

	trn := s.BeginTransaction()

	id, err := trn.MatchCreate(protocol.Match{})
	require.NoError(t, err)

	_, err = trn.MatchRead(id)
	require.NoError(t, err)

	require.False(t, trn.Closed())
	trn.Abort()
	require.Panics(t, func() { trn.Commit() })
	require.Panics(t, func() { trn.Abort() })

	require.True(t, trn.Closed())

	trn = s.BeginTransaction()

	_, err = trn.MatchRead(id)
	require.IsType(t, store.DontExistError(""), err)
}
