package inmemory

import (
	"plateau/store"

	"github.com/sirupsen/logrus"
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
func (s *Transaction) Commit() error {
	if s.Closed() {
		logrus.Panic("You cannot commit a closed transaction")
	}

	*s.inMemory = *s.inMemoryCopy

	s.commitCb(s)
	s.close()

	return nil
}

// Abort implements the `store.Transaction` interface.
func (s *Transaction) Abort() error {
	if s.Closed() {
		logrus.Panic("You cannot abort a closed transaction")
	}

	s.close()

	return nil
}

// Closed implements the `store.Transaction` interface.
func (s *Transaction) Closed() bool {
	return s.closed
}

// Errors implements the `store.Transaction` interface.
func (s *Transaction) Errors() []error {
	return s.errors
}
