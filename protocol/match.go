package protocol

import (
	"math/rand"
	"time"

	"github.com/thoas/go-funk"
)

// Match ...
type Match struct {
	ID string `json:"id"`

	CreatedAt time.Time  `json:"created_at"`
	EndedAt   *time.Time `json:"ended_at"`

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

// IsEnded ...
func (s *Match) IsEnded() bool {
	return s.EndedAt != nil
}

// NextPlayer ...
func (s *Match) NextPlayer(p Player) *Player {
	i := funk.IndexOf(s.Players, p)
	if i == -1 {
		return nil
	}

	if i == len(s.Players)-1 {
		return &s.Players[0]
	}

	return &s.Players[i+1]
}

// RandomPlayer ...
func (s *Match) RandomPlayer() *Player {
	rand.Seed(time.Now().Unix())

	return &s.Players[rand.Intn(len(s.Players))]
}
