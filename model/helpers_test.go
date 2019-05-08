package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsDuplicateError(t *testing.T) {
	t.Parallel()

	require.False(t, IsDuplicateError(errors.New("pq: foobar")))
	require.True(t, IsDuplicateError(errors.New(`pq: duplicate key value violates unique constraint "players_pkey"`)))
}
