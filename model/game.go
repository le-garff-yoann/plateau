package model

import "time"

// Game ...
type Game struct {
	ID uint `json:"id"`

	CreatedAt *time.Time `json:"created_at"`
	EndedAt   *time.Time `json:"ended_at"`

	NumberOfPlayersRequired uint      `json:"number_of_players_required"`
	Players                 []*Player `gorm:"many2many:player_games;" json:"-"`

	Running bool `json:"running"`

	EventContainers []*EventContainer `json:"-"`
}
