package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessageString(t *testing.T) {
	msg := MPlayerAccepts

	require.Equal(t, string(msg), msg.String())
}
