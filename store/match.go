package store

import (
	"time"
)

// Match ...
type Match struct {
	ID string `json:"id"`

	CreatedAt *time.Time `json:"created_at"`
	EndedAt   *time.Time `json:"ended_at"`

	NumberOfPlayersRequired uint      `json:"number_of_players_required"`
	Players                 []*Player `json:"-"`

	Running bool `json:"running"`

	EventContainers []*EventContainer `json:"-"`
}

func (s *Match) String() string {
	return s.ID
}

// MatchStore ...
type MatchStore interface {
	List() (IDs []string, err error)
	Create(Match) (id string, err error)
	Read(matchID string) (*Match, error)

	EndedAt(matchID string, value *time.Time) error

	AddPlayer(matchID string, playerName string) error
	RemovePlayer(matchID string, playerName string) error

	Running(matchID string, value bool) error

	CreateEventContainer(matchID string, ec EventContainer) error
	CreateEventContainerBroadcaster(matchID string) (*EventContainerBroadcaster, error)
}
