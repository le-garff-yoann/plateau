package server

import (
	"plateau/protocol"
)

// NewContext ...
func NewContext() *Context {
	return (&Context{make(map[protocol.Request]func(*MatchRuntime, *protocol.RequestContainer) *protocol.ResponseContainer)})
}

// Context ...
type Context struct {
	handlers map[protocol.Request]func(*MatchRuntime, *protocol.RequestContainer) *protocol.ResponseContainer
}

// Requests ...
func (s *Context) Requests() (requests []protocol.Request) {
	for r := range s.handlers {
		requests = append(requests, r)
	}

	return requests
}

// On ...
func (s *Context) On(request protocol.Request, handlerFunc func(*MatchRuntime, *protocol.RequestContainer) *protocol.ResponseContainer) *Context {
	s.handlers[request] = handlerFunc

	return s
}

// Delete ...
func (s *Context) Delete(requests ...protocol.Request) *Context {
	for _, r := range requests {
		delete(s.handlers, r)
	}

	return s
}

// Complete ...
func (s *Context) Complete(context *Context) *Context {
	for r, handlerFunc := range context.handlers {
		_, ok := s.handlers[r]
		if !ok {
			s.On(r, handlerFunc)
		}
	}

	return s
}

func (s *Context) handle(matchRuntime *MatchRuntime, reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
	handlerFunc, ok := s.handlers[reqContainer.Request]
	if ok {
		return handlerFunc(matchRuntime, reqContainer)
	}

	return &protocol.ResponseContainer{Response: protocol.ResNotImplemented}
}
