package store

import "plateau/protocol"

// DealsChange represents the difference between
// the old and the new state of the same `protocol.Deal`.
type DealsChange struct {
	Old *protocol.Deal `json:"old"`
	New *protocol.Deal `json:"new"`
}

// DealsChangeIterator represents the iterator in its most classical form.
// He is specialized to return only the `protocol.Deal`.
//	- `Next()` fetches the next deal.
//	- `Close()` closes the iterator and stops calls to `Next()`.
type DealsChangeIterator interface {
	Next(*DealsChange) bool
	Close() error
}
