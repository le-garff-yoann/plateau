package server

import (
	"plateau/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContext(t *testing.T) {
	t.Parallel()

	var (
		ctxBeforeHit = 0
		ctxAfterHit  = 0

		ctx = NewContext().
			Before(func(*Context, *protocol.RequestContainer) *protocol.ResponseContainer {
				ctxBeforeHit++

				return nil
			}).
			After(func(*protocol.RequestContainer) *protocol.ResponseContainer {
				ctxAfterHit++

				return nil
			})

		completeCtx = NewContext().On(
			protocol.ReqListRequests,
			func(*protocol.RequestContainer) *protocol.ResponseContainer {
				return &protocol.ResponseContainer{Response: protocol.ResOK}
			},
		)
	)

	require.Empty(t, ctx.handlers)
	require.Len(t, ctx.handlers, len(ctx.Requests()))

	require.Nil(t, ctx.handle(&protocol.RequestContainer{Request: protocol.ReqListRequests}))
	require.NotZero(t, ctxBeforeHit)
	require.Zero(t, ctxAfterHit)

	ctx.OnNotImplemented(func(*protocol.RequestContainer) *protocol.ResponseContainer {
		return &protocol.ResponseContainer{Response: protocol.ResBadRequest}
	})
	require.Equal(
		t,
		ctx.handle(&protocol.RequestContainer{Request: protocol.ReqListRequests}).Response,
		protocol.ResBadRequest,
	)

	ctx.Complete(completeCtx)
	require.Len(t, ctx.handlers, 1)
	require.Len(t, ctx.handlers, len(ctx.Requests()))
	require.Equal(
		t,
		protocol.ResOK,
		ctx.handle(&protocol.RequestContainer{Request: protocol.ReqListRequests}).Response,
	)
	require.NotZero(t, ctxAfterHit)
	require.Equal(
		t,
		protocol.ResBadRequest,
		ctx.handle(&protocol.RequestContainer{Request: protocol.ReqPlayerAccepts}).Response,
	)

	completeCtx.OnNotImplemented(func(*protocol.RequestContainer) *protocol.ResponseContainer {
		return &protocol.ResponseContainer{Response: protocol.ResInternalError}
	})
	ctx.Complete(completeCtx)
	require.Equal(
		t,
		protocol.ResInternalError,
		ctx.handle(&protocol.RequestContainer{Request: protocol.ReqPlayerAccepts}).Response,
	)
}
