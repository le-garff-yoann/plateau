package protocol

import funk "github.com/thoas/go-funk"

// Message ...
type Message struct {
	Code MessageCode `json:"code,omitempty"`

	Payload interface{} `json:"payload,omitempty"`
}

// Concealed ...
func (s *Message) Concealed(playerName ...string) *Message {
	concealedPayload, ok := s.Payload.(MessageConcealedPayload)
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

// MessageCode ...
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
	MPlayerWantToStartTheGame MessageCode = "PLAYER_WANT_TO_START_THE_GAME"
)

func (s MessageCode) String() string {
	return string(s)
}

// MessageConcealedPayload ...
type MessageConcealedPayload struct {
	AllowedNamesCode, AllowedNamesPayload []string

	Data interface{}
}
