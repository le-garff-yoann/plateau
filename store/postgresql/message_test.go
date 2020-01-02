package postgresql

import (
	"plateau/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessageConversion(t *testing.T) {
	t.Parallel()

	msg := &Message{}

	require.IsType(t, &protocol.Message{}, msg.toProtocolStruct())
	require.IsType(t, msg, messageFromProtocolStruct(msg.toProtocolStruct()))
}
