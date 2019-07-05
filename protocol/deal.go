package protocol

import (
	"github.com/thoas/go-funk"
)

// Deal represents an exchange between one or more `Player`.
// The exchange is modeled by a series of `Message`.
//
// *Holder* marks the expectation of a specific player's action.
type Deal struct {
	Holder Player `json:"holder"`

	Messages []Message `json:"messages"`
}

// Find searches for and returns the first `Message` validating *cb*.
func (s *Deal) Find(cb func(Message) bool) *Message {
	f := funk.Find(s.Messages, cb)
	if f == nil {
		return nil
	}

	msg := f.(Message)
	return &msg
}

// FindAll searches for and returns all `Message` validating *cb*.
func (s *Deal) FindAll(cb func(Message) bool) []Message {
	return funk.Filter(s.Messages, cb).([]Message)
}

// FindByMessageCode searches for and returns the first
// `Message` matching *messageCode*.
func (s *Deal) FindByMessageCode(messageCode MessageCode) *Message {
	return s.Find(func(msg Message) bool {
		return msg.Code == messageCode
	})
}

// FindAllByMessageCode searches for and returns the
// all `Message` matching *messageCode*.
func (s *Deal) FindAllByMessageCode(messageCode MessageCode) []Message {
	return s.FindAll(func(msg Message) bool {
		return msg.Code == messageCode
	})
}

// IsActive returns `true` is the deal is "finished".
func (s *Deal) IsActive() bool {
	return s.FindByMessageCode(MDealCompleted) == nil && s.FindByMessageCode(MDealAborded) == nil
}

// WithMessagesConcealed returns itself with his `Message` concealed.
func (s *Deal) WithMessagesConcealed(playerName ...string) *Deal {
	deal := Deal{
		Holder:   s.Holder,
		Messages: []Message{},
	}

	for _, msg := range s.Messages {
		deal.Messages = append(deal.Messages, *msg.Concealed(playerName...))
	}

	return &deal
}

// IndexDeals indexes a collection of `Deal`.
func IndexDeals(deals []Deal, i uint) *Deal {
	i++

	if int(i) > len(deals) {
		return nil
	}

	return &deals[len(deals)-int(i)]
}
