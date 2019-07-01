package store

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDuplicateError(t *testing.T) {
	t.Parallel()

	require.Error(t, DuplicateError(""))
}

func TestDontExistError(t *testing.T) {
	t.Parallel()

	require.Error(t, DontExistError(""))
}

func TestPlayerParticipationError(t *testing.T) {
	t.Parallel()

	require.Error(t, PlayerParticipationError(""))
}
