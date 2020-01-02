package postgresql

import (
	"os"
	"plateau/store"
	"testing"

	"github.com/stretchr/testify/require"
)

func testStr(t *testing.T) *store.TestStore {
	t.Parallel()

	postgresqlURL = os.Getenv("TEST_PG_DSN")
	if len(postgresqlURL) == 0 {
		t.Skip("Assign a PostgreSQL connection string to $TEST_PG_DSN")
	}

	createTableOptionsTemp = true

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
