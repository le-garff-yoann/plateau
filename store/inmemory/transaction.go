package inmemory

import (
	"plateau/store"
)

// Transaction implements the `store.Transaction` interface.
type Transaction struct {
	errors []error

	inMemory, inMemoryCopy *inMemory
	dealChangeSubmitter    func(*store.DealsChange)

	closed bool
	done   func()
}

func (s *Transaction) close() {
	s.closed = true

	s.done()
}

// Commit implements the `store.Transaction` interface.
func (s *Transaction) Commit() {
	if s.Closed() {
		panic("You cannot commit a closed transaction")
	}

	*s.inMemory = *s.inMemoryCopy
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
