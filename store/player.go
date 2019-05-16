package store

// Player ...
type Player struct {
	Name string `json:"name"`

	Password string `json:"-"`

	Wins  uint `json:"wins"`
	Loses uint `json:"loses"`
	Ties  uint `json:"ties"`
}

func (s *Player) String() string {
	return s.Name
}

// PlayerStore ...
type PlayerStore interface {
	List() (names []string, errs error)
	Create(Player) error
	Read(playerName string) (*Player, error)

	IncreaseWins(playerName string, increase uint) error
	IncreaseLoses(playerName string, increase uint) error
	IncreaseTies(playerName string, increase uint) error
}
