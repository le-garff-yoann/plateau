package inmemory

import (
	"plateau/store"
)

// Transaction ...
type Transaction struct {
	errors []error

	inMemory, inMemoryCopy *inMemory
	dealChangeSubmitter    func(*store.DealChange)

	closed bool
	done   func()
}

func (s *Transaction) close() {
	s.closed = true

	s.done()
}

// Commit ...
func (s *Transaction) Commit() {
	if s.Closed() {
		panic("You cannot commit a closed transaction")
	}

	*s.inMemory = *s.inMemoryCopy
	s.close()
}

// Abort ...
func (s *Transaction) Abort() {
	if s.Closed() {
		panic("You cannot abort a closed transaction")
	}

	s.close()
}

// Closed ...
func (s *Transaction) Closed() bool {
	return s.closed
}

// Errors ...
func (s *Transaction) Errors() []error {
	return s.errors
}
