package model

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm/dialects/postgres"
)

// Event ...
type Event string

const (
	// EInternalError ...
	EInternalError Event = "INTERNAL_ERROR"
	// EProcessed ...
	EProcessed Event = "PROCESSED"
	// EIllegal ...
	EIllegal Event = "ILLEGAL"
	// EListEvents ...
	EListEvents Event = "LIST"
	// EGameWantToStart ...
	EGameWantToStart Event = "GAME_WANT_TO_START"
	// EGameStarts ...
	EGameStarts Event = "GAME_STARTS"
	// EGameEnds ...
	EGameEnds Event = "GAME_ENDS"
	// EPlayerConnects ...
	EPlayerConnects Event = "PLAYER_CONNECTS"
	// EPlayerDisconnects ...
	EPlayerDisconnects Event = "PLAYER_DISCONNECTS"
	// EPlayerReconnects ...
	EPlayerReconnects Event = "PLAYER_RECONNECTS"
	// EPlayerWantToJoin ...
	EPlayerWantToJoin Event = "PLAYER_WANT_TO_JOIN"
	// EPlayerJoins ...
	EPlayerJoins Event = "PLAYER_JOINS"
	// EPlayerWantToLeave ...
	EPlayerWantToLeave Event = "PLAYER_WANT_TO_LEAVE"
	// EPlayerLeaves ...
	EPlayerLeaves Event = "PLAYER_LEAVES"
	// EPlayerWantToSurrender ...
	EPlayerWantToSurrender Event = "PLAYER_WANT_TO_SURRENDER"
	// EPlayerSurrenders ...
	EPlayerSurrenders Event = "PLAYER_SURRENDERS"
)

func (s Event) String() string { return string(s) }

// EventContainer ...
type EventContainer struct {
	ID uint `json:"-"`

	Event `json:"event"`

	GameID uint  `json:"-"`
	Game   *Game `json:"-"`

	Emitter   *Player   `json:"emitter,omitempty"`
	Receivers []*Player `json:"receivers,omitempty"`
	Subjects  []*Player `json:"subjects,omitempty"`

	Payload postgres.Jsonb `json:"payload,omitempty"`
}

func (s EventContainer) String() string {
	r := []string{fmt.Sprintf(`Event: "%s"`, s.Event)}

	if s.Emitter.Name != "" {
		r = append(r, fmt.Sprintf(`Emitter: "%s"`, s.Emitter.Name))
	}

	if len(s.Receivers) > 0 {
		var names []string
		for _, p := range s.Receivers {
			names = append(names, p.Name)
		}

		r = append(r, fmt.Sprintf(`Receivers: "%s"`, strings.Join(names, ", ")))
	}

	if len(s.Subjects) > 0 {
		var names []string
		for _, p := range s.Receivers {
			names = append(names, p.Name)
		}

		r = append(r, fmt.Sprintf(`Subjects: "%s"`, strings.Join(names, ", ")))
	}

	return strings.Join(r, " - ")
}

// IsLegal ...
func (s EventContainer) IsLegal() bool {
	switch s.Event {
	case
		EIllegal, EInternalError,
		EGameWantToStart, EGameStarts, EGameEnds,
		EPlayerConnects, EPlayerDisconnects, EPlayerReconnects,
		EPlayerJoins, EPlayerLeaves, EPlayerSurrenders:
		return s.Emitter == nil
	case EListEvents, EPlayerWantToJoin, EPlayerWantToLeave, EPlayerWantToSurrender:
		return true
	}

	return false
}
