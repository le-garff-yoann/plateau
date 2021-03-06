package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestString(t *testing.T) {
	t.Parallel()

	req := ReqListRequests

	require.Equal(t, string(req), req.String())
}

func TestRequestContainerString(t *testing.T) {
	t.Parallel()

	reqContainer := RequestContainer{Request: ReqListRequests}

	require.Equal(t, string(reqContainer.Request), reqContainer.String())
}
