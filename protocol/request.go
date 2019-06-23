package protocol

// Request ...
type Request string

const (
	// ReqPlayerAccepts ...
	ReqPlayerAccepts Request = "PLAYER_ACCEPTS"
	// ReqPlayerRefuses ...
	ReqPlayerRefuses Request = "PLAYER_REFUSES"
	// ReqListRequests ...
	ReqListRequests Request = "LIST_REQUESTS"
	// ReqGetCurrentTransaction ...
	ReqGetCurrentTransaction Request = "GET_CURRENT_TRANSACTION"
	// ReqPlayerWantToJoin ...
	ReqPlayerWantToJoin Request = "PLAYER_WANT_TO_JOIN"
	// ReqPlayerWantToLeave ...
	ReqPlayerWantToLeave Request = "PLAYER_WANT_TO_LEAVE"
	// ReqPlayerWantToStartTheGame ...
	ReqPlayerWantToStartTheGame Request = "PLAYER_WANT_TO_START_THE_GAME"
)

func (s Request) String() string {
	return string(s)
}

// RequestContainer ...
type RequestContainer struct {
	// ID      uuid.UUID `json:"id"`
	Request `json:"request"`

	Player *Player `json:"-"`
	Match  *Match  `json:"-"`
}

func (s RequestContainer) String() string {
	return string(s.Request)
}
