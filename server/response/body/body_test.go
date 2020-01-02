package body

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResponse(t *testing.T) {
	t.Parallel()

	var (
		body = New()

		ok = ""
		ko = errors.New("")
	)

	require.Empty(t, body.Ok().Successes)
	require.Len(t, body.Ok(ok).Successes, 1)
	require.Len(t, body.Ok(ok).Successes, 2)
	require.Len(t, body.Ok(ok, ok).Successes, 4)

	require.Empty(t, body.Ko().Failures)
	require.Len(t, body.Ko(ko).Failures, 1)
	require.Len(t, body.Ko(ko).Failures, 2)
	require.Len(t, body.Ko(ko, ko, nil).Failures, 4)
}
