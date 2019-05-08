package model

// Player ...
type Player struct {
	Name     string `gorm:"primary_key" gorm:"unique" json:"name"`
	Password string `json:"-"`

	Wins  uint `json:"wins"`
	Loses uint `json:"loses"`
	Ties  uint `json:"ties"`

	Games []*Game `gorm:"many2many:player_games;" json:"-"`
}
