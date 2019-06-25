package protocol

import (
	"time"
)

// Match ...
type Match struct {
	ID string `json:"id"`

	CreatedAt time.Time  `json:"created_at"`
	EndedAt   *time.Time `json:"ended_at"`

	ConnectedPlayers []Player `json:"-"`

	NumberOfPlayersRequired uint     `json:"number_of_players_required"`
	Players                 []Player `json:"-"`

	Deals []Deal `json:"-"`
}

func (s *Match) String() string {
	return s.ID
}

// IsFull ...
func (s *Match) IsFull() bool {
	return int(s.NumberOfPlayersRequired) == len(s.Players)
}
