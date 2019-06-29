package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestString(t *testing.T) {
	req := ReqListRequests

	require.Equal(t, string(req), req.String())
}

func TestRequestContainerString(t *testing.T) {
	reqContainer := RequestContainer{Request: ReqListRequests}

	require.Equal(t, string(reqContainer.Request), reqContainer.String())
}
