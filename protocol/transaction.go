package protocol

import (
	"github.com/thoas/go-funk"
)

// Transaction ...
type Transaction struct {
	Holder Player `json:"holder"`

	Messages []Message `json:"messages"`
}

// Find ...
func (s *Transaction) Find(cb func(Message) bool) *Message {
	f := funk.Find(s.Messages, cb)
	if f == nil {
		return nil
	}

	msg := f.(Message)
	return &msg
}

// FindAll ...
func (s *Transaction) FindAll(cb func(Message) bool) []Message {
	return funk.Filter(s.Messages, cb).([]Message)
}

// FindByMessageCode ...
func (s *Transaction) FindByMessageCode(messageCode MessageCode) *Message {
	return s.Find(func(msg Message) bool {
		return msg.MessageCode == messageCode
	})
}

// FindAllByMessageCode ...
func (s *Transaction) FindAllByMessageCode(messageCode MessageCode) []Message {
	return s.FindAll(func(msg Message) bool {
		return msg.MessageCode == messageCode
	})
}

// IsActive ...
func (s *Transaction) IsActive() bool {
	return s.FindByMessageCode(MTransactionCompleted) == nil && s.FindByMessageCode(MTransactionAborded) == nil
}

// IndexTransactions ...
func IndexTransactions(transactions []Transaction, i uint) *Transaction {
	i++

	if int(i) > len(transactions) {
		return nil
	}

	return &transactions[len(transactions)-int(i)]
}
