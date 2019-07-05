package protocol

import funk "github.com/thoas/go-funk"

// Message represents a change, an action
// or other within a `Deal`.
type Message struct {
	Code MessageCode `json:"code,omitempty"`

	Payload interface{} `json:"payload,omitempty"`
}

// Concealed convert itself into a "normal" message if the *Payload*
// is assertable to `ConcealedMessagePayload`.
func (s *Message) Concealed(playerName ...string) *Message {
	concealedPayload, ok := s.Payload.(ConcealedMessagePayload)
	if ok {
		var (
			msg = Message{}

			allowed = func(names []string) bool {
				return len(playerName) == 0 || names == nil || funk.Find(names, func(name string) bool {
					return name == playerName[0]
				}) != nil
			}
		)

		if allowed(concealedPayload.AllowedNamesCode) {
			msg.Code = s.Code
		}

		if allowed(concealedPayload.AllowedNamesPayload) {
			msg.Payload = concealedPayload.Data
		}

		return &msg
	}

	return s
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

// ConcealedMessagePayload allows to display parts of a `Message`
// only to a specific list of `Player`.
//	- *AllowedNamesCode* is the list of players allowed to see
//	the `Message.MessageCode`.
// 	- *AllowedNamesPayload* is the list of players allowed to see
//	the `Message.Payload`.
//	- *Data* is analogous to the `Message.Payload`
//	of a "normal" message.
type ConcealedMessagePayload struct {
	AllowedNamesCode, AllowedNamesPayload []string

	Data interface{}
}
