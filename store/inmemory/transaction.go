package inmemory

import (
	"plateau/store"
)

// Transaction implements the `store.Transaction` interface.
type Transaction struct {
	errors []error

	matchNotifications []store.MatchNotification

	inMemory, inMemoryCopy *inMemory

	closed         bool
	commitCb, done func(*Transaction)
}

func (s *Transaction) close() {
	s.closed = true

	s.done(s)
}

// Commit implements the `store.Transaction` interface.
func (s *Transaction) Commit() {
	if s.Closed() {
		panic("You cannot commit a closed transaction")
	}

	*s.inMemory = *s.inMemoryCopy

	s.commitCb(s)
	s.close()
}

// Abort implements the `store.Transaction` interface.
func (s *Transaction) Abort() {
	if s.Closed() {
		panic("You cannot abort a closed transaction")
	}

	s.close()
}

// Closed implements the `store.Transaction` interface.
func (s *Transaction) Closed() bool {
	return s.closed
}

// Errors implements the `store.Transaction` interface.
func (s *Transaction) Errors() []error {
	return s.errors
}
