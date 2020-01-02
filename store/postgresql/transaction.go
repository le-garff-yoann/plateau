package postgresql

import (
	"plateau/store"

	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
)

// Transaction implements the `store.Transaction` interface.
type Transaction struct {
	errors []error

	matchNotifications []store.MatchNotification

	tx *pg.Tx

	closed   bool
	commitCb func(*Transaction) error
}

func (s *Transaction) close() {
	s.closed = true
}

// Commit implements the `store.Transaction` interface.
func (s *Transaction) Commit() error {
	if s.Closed() {
		logrus.Panic("You cannot commit a closed transaction")
	}

	if err := s.tx.Commit(); err != nil {
		return err
	}

	if err := s.commitCb(s); err != nil {
		return err
	}
	s.close()

	return nil
}

// Abort implements the `store.Transaction` interface.
func (s *Transaction) Abort() error {
	if s.Closed() {
		logrus.Panic("You cannot abort a closed transaction")
	}

	if err := s.tx.Rollback(); err != nil {
		return err
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
