package inmemory

import (
	"plateau/store"
	"testing"

	"github.com/stretchr/testify/require"
)

func testStr(t *testing.T) *store.TestStore {
	t.Parallel()

	return &store.TestStore{T: t, Str: &Store{}}
}

func TestStore(t *testing.T) {
	testStr := testStr(t)

	require.NoError(t, testStr.Str.Open())
	require.NoError(t, testStr.Str.Close())
}

func TestNotificationsChannel(t *testing.T) {
	testStr(t).TestNotificationsChannel()
}
