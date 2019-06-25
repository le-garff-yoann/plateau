package protocol

// Message ...
type Message struct {
	MessageCode `json:"message_code"`

	Payload interface{} `json:"payload"`
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
