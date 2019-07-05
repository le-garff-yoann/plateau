package rockpaperscissors

import "plateau/protocol"

const (
	// ReqRock is the request that generates `MRock`.
	ReqRock protocol.Request = "PLAY_ROCK"
	// ReqPaper is the request that generates `MPaper`.
	ReqPaper protocol.Request = "PLAY_PAPER"
	// ReqScissors is the request that generates `MScissors`.
	ReqScissors protocol.Request = "PLAY_SCISSORS"

	// MRock represents a "rock" played.
	MRock protocol.MessageCode = "ROCK"
	// MPaper represents a "paper" played.
	MPaper protocol.MessageCode = "PAPER"
	// MScissors represents a "paper" played.
	MScissors protocol.MessageCode = "SCISSORS"
)
