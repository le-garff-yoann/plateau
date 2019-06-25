package store

import "plateau/protocol"

// DealChange ...
type DealChange struct {
	Old *protocol.Deal
	New *protocol.Deal
}

// NewMessages ...
func (s *DealChange) NewMessages() (msgs []*protocol.Message) {
	for i := len(s.Old.Messages); i < len(s.New.Messages); i++ {
		msgs = append(msgs, &s.New.Messages[i])
	}

	return msgs
}

// DealChangeIterator ...
type DealChangeIterator interface {
	Next(*DealChange) bool
	Close() error
}
