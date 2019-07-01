package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessageConcealed(t *testing.T) {
	t.Parallel()

	var (
		msg = Message{Code: MPlayerAccepts, Payload: "foo"}

		playerAName = "bar"
		playerBName = "baz"
	)

	require.Equal(t, &msg, msg.Concealed(playerAName))

	concealedMsg := Message{
		Code: msg.Code,
		Payload: MessageConcealedPayload{
			AllowedNamesCode:    nil,
			AllowedNamesPayload: nil,
			Data:                msg.Payload.(string),
		},
	}
	require.Equal(t, &msg, concealedMsg.Concealed(playerAName))

	mutateConcealedMsgPayload := func(mutater func(*MessageConcealedPayload)) {
		p := concealedMsg.Payload.(MessageConcealedPayload)

		mutater(&p)

		concealedMsg.Payload = p
	}

	require.Equal(t, msg.Code, concealedMsg.Concealed().Code)

	mutateConcealedMsgPayload(func(m *MessageConcealedPayload) {
		m.AllowedNamesCode = []string{playerBName}
	})
	require.Empty(t, concealedMsg.Concealed(playerAName).Code)
	require.Equal(t, msg.Code, concealedMsg.Concealed().Code)

	mutateConcealedMsgPayload(func(m *MessageConcealedPayload) {
		m.AllowedNamesCode = []string{playerAName}
	})
	require.Equal(t, msg.Code, concealedMsg.Concealed(playerAName).Code)

	mutateConcealedMsgPayload(func(m *MessageConcealedPayload) {
		m.AllowedNamesPayload = []string{playerBName}
	})
	require.Empty(t, concealedMsg.Concealed(playerAName).Payload)

	mutateConcealedMsgPayload(func(m *MessageConcealedPayload) {
		m.AllowedNamesPayload = []string{playerAName}
	})
	require.Equal(t, msg.Payload, concealedMsg.Concealed(playerAName).Payload)
}

func TestMessageCodeString(t *testing.T) {
	t.Parallel()

	msgCode := MPlayerAccepts

	require.Equal(t, string(msgCode), msgCode.String())
}
