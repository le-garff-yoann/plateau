package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessageString(t *testing.T) {
	t.Parallel()

	msg := MPlayerAccepts

	require.Equal(t, string(msg), msg.String())
}
