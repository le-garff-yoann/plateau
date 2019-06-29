package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResponseString(t *testing.T) {
	res := ResOK

	require.Equal(t, string(res), res.String())
}

func TestResponseContainerString(t *testing.T) {
	resContainer := ResponseContainer{Response: ResOK}

	require.Equal(t, string(resContainer.Response), resContainer.String())
}
