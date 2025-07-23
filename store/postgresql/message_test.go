package postgresql

import (
	"testing"

	"github.com/le-garff-yoann/plateau/protocol"
	"github.com/stretchr/testify/require"
)

func TestMessageConversion(t *testing.T) {
	t.Parallel()

	msg := &Message{}

	require.IsType(t, &protocol.Message{}, msg.toProtocolStruct())
	require.IsType(t, msg, messageFromProtocolStruct(msg.toProtocolStruct()))
}
