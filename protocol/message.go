package protocol

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// Message represents a change, an action
// or other within a `Deal`.
//	- *AllowedNamesCode* is the list of players allowed to see
//	the `Code`.
// 	- *AllowedNamesPayload* is the list of players allowed to see
//	the `Payload`.
type Message struct {
	Code    MessageCode     `json:"code,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`

	AllowedNamesCode    []string `json:"-"`
	AllowedNamesPayload []string `json:"-"`
}

// DecodePayload writes the decoded value of *Payload* to **p*.
func (s *Message) DecodePayload(p interface{}) {
	if err := json.Unmarshal(s.Payload, p); err != nil {
		logrus.Fatal(err)
	}
}

// EncodePayload set *Payload* to the encoded value of *v*.
func (s *Message) EncodePayload(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		logrus.Fatal(err)
	}

	s.Payload = b
}

// Concealed convert itself into a "normal" message if the *Payload*
// is "assertable" to `ConcealedMessagePayload`.
func (s *Message) Concealed(playerName ...string) *Message {
	var (
		msg = Message{}

		allowed = func(names []string) bool {
			if len(playerName) > 0 && len(names) > 0 {
				for _, name := range names {
					if name == playerName[0] {
						return true
					}
				}

				return false
			}

			return true
		}
	)

	if allowed(s.AllowedNamesCode) {
		msg.Code = s.Code
	}

	if allowed(s.AllowedNamesPayload) {
		msg.Payload = s.Payload
	}

	return &msg
}

// MessageCode is somehow the identifier of `Message`.
// It helps to identify its nature.
type MessageCode string

const (
	// MDealCompleted ...
	MDealCompleted MessageCode = "DEAL_COMPLETED"
	// MDealAborded ...
	MDealAborded MessageCode = "DEAL_ABORDED"
	// MPlayerAccepts ...
	MPlayerAccepts MessageCode = "PLAYER_ACCEPTS"
	// MPlayerRefuses ...
	MPlayerRefuses MessageCode = "PLAYER_REFUSES"
	// MPlayerWantToStartTheGame ...
	MPlayerWantToStartTheGame MessageCode = "PLAYER_WANT_TO_START_THE_MATCH"
)

func (s MessageCode) String() string {
	return string(s)
}
