package protocol

// Request represents a request emitted by a `Player`.
type Request string

const (
	// ReqListRequests means that a `Player` wants
	// to display the list of available `Request`.
	ReqListRequests Request = "?" // REFACTOR: Must have a dedicated endpoint?
	// ReqPlayerAccepts means that a `Player` accepts something.
	ReqPlayerAccepts Request = "PLAYER_ACCEPTS"
	// ReqPlayerRefuses means that a `Player` refuses something.
	ReqPlayerRefuses Request = "PLAYER_REFUSES"
	// ReqPlayerWantToJoin means that a `Player` want to join a `Match`.
	ReqPlayerWantToJoin Request = "PLAYER_WANT_TO_JOIN"
	// ReqPlayerWantToLeave means that a `Player` want to leave a `Match`.
	ReqPlayerWantToLeave Request = "PLAYER_WANT_TO_LEAVE"
	// ReqPlayerWantToStartTheMatch means that a `Player` want to start a `Match`.
	ReqPlayerWantToStartTheMatch Request = "PLAYER_WANT_TO_START_THE_MATCH"
)

func (s Request) String() string {
	return string(s)
}

// RequestContainer is `Request`'s container.
type RequestContainer struct {
	Request `json:"request"`

	Player *Player `json:"-"`
	Match  *Match  `json:"-"`
}

func (s RequestContainer) String() string {
	return string(s.Request)
}
