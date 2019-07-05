package store

import "plateau/protocol"

// DealChange represents the difference between
// the old and the new state of the same `protocol.Deal`.
type DealChange struct {
	Old *protocol.Deal
	New *protocol.Deal
}

// DealChangeIterator represents the iterator in its most classical form.
// He is specialized to return only the `protocol.Deal`.
//	- `Next()` fetches the next deal.
//	- `Close()` closes the iterator and stops calls to `Next()`.
type DealChangeIterator interface {
	Next(*DealChange) bool
	Close() error
}
