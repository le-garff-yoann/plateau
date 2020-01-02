package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessageCodeString(t *testing.T) {
	t.Parallel()

	msgCode := MPlayerAccepts

	require.Equal(t, string(msgCode), msgCode.String())
}

func TestMessagePayloadEncodingAndDecoding(t *testing.T) {
	t.Parallel()

	var (
		payload = "foo"

		msg = Message{Code: MPlayerAccepts}
	)

	require.NotPanics(t, func() { msg.EncodePayload(payload) })

	var decodedPayload string
	require.NotPanics(t, func() { msg.DecodePayload(&decodedPayload) })
	require.Equal(t, payload, decodedPayload)
}

func TestMessageConcealed(t *testing.T) {
	t.Parallel()

	var (
		msg = Message{Code: MPlayerAccepts}

		playerAName = "bar"
		playerBName = "baz"
	)

	newMsg := func(mutater func(*Message)) {
		mutater(&Message{
			Code:    msg.Code,
			Payload: msg.Payload,
		})
	}

	newMsg(func(m *Message) {
		require.Equal(t, &msg, m.Concealed(playerAName))
	})

	newMsg(func(m *Message) {
		m.AllowedNamesCode = append(m.AllowedNamesCode, playerAName)

		require.Equal(t, &msg, m.Concealed(playerAName))
		require.Empty(t, m.Concealed(playerBName).Code)
		require.Equal(t, msg.Payload, m.Concealed(playerBName).Payload)

		m.AllowedNamesCode = append(m.AllowedNamesCode, playerBName)

		require.Equal(t, &msg, m.Concealed(playerAName))
		require.Equal(t, &msg, m.Concealed(playerBName))
	})

	newMsg(func(m *Message) {
		m.AllowedNamesPayload = append(m.AllowedNamesPayload, playerAName)

		require.Equal(t, &msg, m.Concealed(playerAName))
		require.Empty(t, m.Concealed(playerBName).Payload)
		require.Equal(t, msg.Code, m.Concealed(playerBName).Code)

		m.AllowedNamesPayload = append(m.AllowedNamesPayload, playerBName)

		require.Equal(t, &msg, m.Concealed(playerAName))
		require.Equal(t, &msg, m.Concealed(playerBName))
	})
}
