package protocol

// Player ...
type Player struct {
	Name     string `json:"name"`
	Password string `json:"-"`

	Wins  uint `json:"wins"`
	Loses uint `json:"loses"`
	Ties  uint `json:"ties"`
}

func (s *Player) String() string {
	return s.Name
}
