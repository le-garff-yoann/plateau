package protocol

import (
	"github.com/thoas/go-funk"
)

// Deal ...
type Deal struct {
	Holder Player `json:"holder"`

	Messages []Message `json:"messages"`
}

// Find ...
func (s *Deal) Find(cb func(Message) bool) *Message {
	f := funk.Find(s.Messages, cb)
	if f == nil {
		return nil
	}

	msg := f.(Message)
	return &msg
}

// FindAll ...
func (s *Deal) FindAll(cb func(Message) bool) []Message {
	return funk.Filter(s.Messages, cb).([]Message)
}

// FindByMessageCode ...
func (s *Deal) FindByMessageCode(messageCode MessageCode) *Message {
	return s.Find(func(msg Message) bool {
		return msg.MessageCode == messageCode
	})
}

// FindAllByMessageCode ...
func (s *Deal) FindAllByMessageCode(messageCode MessageCode) []Message {
	return s.FindAll(func(msg Message) bool {
		return msg.MessageCode == messageCode
	})
}

// IsActive ...
func (s *Deal) IsActive() bool {
	return s.FindByMessageCode(MDealCompleted) == nil && s.FindByMessageCode(MDealAborded) == nil
}

// IndexDeals ...
func IndexDeals(deals []Deal, i uint) *Deal {
	i++

	if int(i) > len(deals) {
		return nil
	}

	return &deals[len(deals)-int(i)]
}
