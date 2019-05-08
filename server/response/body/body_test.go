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

	require.Len(t, body.Ok().Successes, 0)
	require.Len(t, body.Ok(ok).Successes, 1)
	require.Len(t, body.Ok(ok).Successes, 2)
	require.Len(t, body.Ok(ok, ok).Successes, 4)

	require.Len(t, body.Ko().Failures, 0)
	require.Len(t, body.Ko(ko).Failures, 1)
	require.Len(t, body.Ko(ko).Failures, 2)
	require.Len(t, body.Ko(ko, ko).Failures, 4)
}
