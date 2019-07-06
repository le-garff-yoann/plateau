package store

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDuplicateError(t *testing.T) {
	t.Parallel()

	require.Error(t, DuplicateError(""))
	require.Equal(t, "foo", string(DuplicateError("foo")))
}

func TestDontExistError(t *testing.T) {
	t.Parallel()

	require.Error(t, DontExistError(""))
	require.Equal(t, "foo", string(DontExistError("foo")))
}

func TestPlayerParticipationError(t *testing.T) {
	t.Parallel()

	require.Error(t, PlayerParticipationError(""))
	require.Equal(t, "foo", string(PlayerParticipationError("foo")))
}
