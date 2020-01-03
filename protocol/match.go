package protocol

import (
	"math/rand"
	"time"
)

// Match is the representation of a match,
// with its beginning, its end, its `Player` and
// its interactions aka `Deal`.
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

// IsFull returns `true` if the number of `Player`
// present equals `Match.NumberOfPlayersRequired`.
func (s *Match) IsFull() bool {
	return int(s.NumberOfPlayersRequired) == len(s.Players)
}

// IsEnded returns `true` if `Match.EndedAt` isn't `nil`.
func (s *Match) IsEnded() bool {
	return s.EndedAt != nil
}

// NextPlayer returns the next players after *p* in
// the `Match.Players` slice.
func (s *Match) NextPlayer(p Player) *Player {
	for i, player := range s.Players {
		if p.Name == player.Name {
			if i == len(s.Players)-1 {
				return &s.Players[0]
			}

			return &s.Players[i+1]
		}
	}

	return nil
}

// RandomPlayer returns a random `Match.Players` from
// the `Match.Players` slice.
func (s *Match) RandomPlayer() *Player {
	rand.Seed(time.Now().Unix())

	return &s.Players[rand.Intn(len(s.Players))]
}
