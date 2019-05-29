package server

import (
	"plateau/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContext(t *testing.T) {
	t.Parallel()

	var (
		mRuntime = MatchRuntime{}

		ctx = NewContext()

		completeCtx = NewContext().On(
			protocol.ReqListRequests,
			func(*MatchRuntime, *protocol.RequestContainer) *protocol.ResponseContainer {
				return &protocol.ResponseContainer{Response: protocol.ResOK}
			},
		)
	)

	require.Empty(t, ctx.handlers)
	require.Len(t, ctx.handlers, len(ctx.Requests()))

	ctx.Complete(completeCtx)
	require.Len(t, ctx.handlers, 1)
	require.Len(t, ctx.handlers, len(ctx.Requests()))
	require.Equal(
		t,
		ctx.handle(&mRuntime, &protocol.RequestContainer{Request: protocol.ReqListRequests}).Response,
		protocol.ResOK,
	)

	require.Empty(t, ctx.Delete(protocol.ReqListRequests).handlers)
}
