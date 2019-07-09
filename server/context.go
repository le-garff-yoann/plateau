package server

import (
	"plateau/protocol"
)

// NewContext returns a new `Context`.
func NewContext() *Context {
	return (&Context{
		nil, nil, nil,
		make(map[protocol.Request]func(*protocol.RequestContainer) *protocol.ResponseContainer),
	})
}

// Context is analogous to a router.
// It embeds one or more handlers that are to be executed according to the `protocol.RequestContainer` received.
type Context struct {
	beforeHandler                       func(*Context, *protocol.RequestContainer) *protocol.ResponseContainer
	afterHandler, notImplementedHandler func(*protocol.RequestContainer) *protocol.ResponseContainer
	handlers                            map[protocol.Request]func(*protocol.RequestContainer) *protocol.ResponseContainer
}

// Requests lists the `protocol.Request` on which the `Context` will execute.
func (s *Context) Requests() (requests []protocol.Request) {
	for r := range s.handlers {
		requests = append(requests, r)
	}

	return requests
}

// Before registers a handler that will execute before all others.
func (s *Context) Before(handlerFunc func(*Context, *protocol.RequestContainer) *protocol.ResponseContainer) *Context {
	s.beforeHandler = handlerFunc

	return s
}

// After registers a handler that will execute after all others.
func (s *Context) After(handlerFunc func(*protocol.RequestContainer) *protocol.ResponseContainer) *Context {
	s.afterHandler = handlerFunc

	return s
}

// OnNotImplemented registers a handler that will execute on an unregistered `protocol.Request`.
func (s *Context) OnNotImplemented(handlerFunc func(*protocol.RequestContainer) *protocol.ResponseContainer) *Context {
	s.notImplementedHandler = handlerFunc

	return s
}

// On registers a handler that will execute on a particular `protocol.Request`.
func (s *Context) On(request protocol.Request, handlerFunc func(*protocol.RequestContainer) *protocol.ResponseContainer) *Context {
	s.handlers[request] = handlerFunc

	return s
}

// Complete merges itself with another `Context`. Priority on merging is given to the new one.
func (s *Context) Complete(ctx *Context) *Context {
	if ctx.beforeHandler != nil {
		s.beforeHandler = ctx.beforeHandler
	}

	if ctx.afterHandler != nil {
		s.afterHandler = ctx.afterHandler
	}

	if ctx.notImplementedHandler != nil {
		s.notImplementedHandler = ctx.notImplementedHandler
	}

	for r, handlerFunc := range ctx.handlers {
		if _, ok := s.handlers[r]; !ok {
			s.On(r, handlerFunc)
		}
	}

	return s
}

func (s *Context) handle(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
	if s.beforeHandler != nil {
		if res := s.beforeHandler(s, reqContainer); res != nil {
			return res
		}
	}

	if handlerFunc, ok := s.handlers[reqContainer.Request]; ok {
		res := handlerFunc(reqContainer)

		if s.afterHandler != nil {
			if afterRes := s.afterHandler(reqContainer); afterRes != nil {
				res = afterRes
			}
		}

		return res
	}

	if s.notImplementedHandler != nil {
		return s.notImplementedHandler(reqContainer)
	}

	return nil
}
