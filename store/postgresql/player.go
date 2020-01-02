package postgresql

import (
	"fmt"
	"plateau/protocol"
	"plateau/store"

	"github.com/go-pg/pg"
)

// Player ...
type Player struct {
	ID       string
	Password string

	Wins  uint
	Loses uint
	Ties  uint
}

func playerFromProtocolStruct(m *protocol.Player) *Player {
	return &Player{
		ID:       m.Name,
		Password: m.Password,
		Wins:     m.Wins,
		Loses:    m.Loses,
		Ties:     m.Ties,
	}
}

func (s *Player) toProtocolStruct() *protocol.Player {
	return &protocol.Player{
		Name:     s.ID,
		Password: s.Password,
		Wins:     s.Wins,
		Loses:    s.Loses,
		Ties:     s.Ties,
	}
}

// PlayerList implements the `store.Transaction` interface.
func (s *Transaction) PlayerList() (names []string, err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}
	}()

	var players []Player
	if err := s.tx.Model(&players).Select(); err != nil {
		return nil, err
	}

	for _, p := range players {
		names = append(names, p.ID)
	}

	return names, err
}

// PlayerCreate implements the `store.Transaction` interface.
func (s *Transaction) PlayerCreate(p protocol.Player) (err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}
	}()

	if err = s.tx.Insert(playerFromProtocolStruct(&p)); err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() {
			return store.DuplicateError(fmt.Sprintf(`The player %s already exist`, p.Name))
		}
	}

	return
}

// PlayerRead implements the `store.Transaction` interface.
func (s *Transaction) PlayerRead(name string) (_ *protocol.Player, err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}
	}()

	player := Player{ID: name}
	if err = s.tx.
		Model(&player).
		WherePK().
		Select(); err != nil {
		if err == pg.ErrNoRows {
			return nil, store.DontExistError(fmt.Sprintf(`The player %s doesn't exist`, name))
		}

		return nil, err
	}

	return player.toProtocolStruct(), nil
}

// PlayerIncreaseWins implements the `store.Transaction` interface.
func (s *Transaction) PlayerIncreaseWins(name string, increase uint) (err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}
	}()

	player := Player{ID: name}
	if err = s.tx.Select(&player); err != nil {
		if err == pg.ErrNoRows {
			return store.DontExistError(fmt.Sprintf(`The player %s doesn't exist`, name))
		}

		return
	}

	player.Wins += increase

	_, err = s.tx.
		Model(&player).
		Column("wins").
		WherePK().
		Update()

	return
}

// PlayerIncreaseLoses implements the `store.Transaction` interface.
func (s *Transaction) PlayerIncreaseLoses(name string, increase uint) (err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}
	}()

	player := Player{ID: name}
	if err = s.tx.Select(&player); err != nil {
		if err == pg.ErrNoRows {
			return store.DontExistError(fmt.Sprintf(`The player %s doesn't exist`, name))
		}

		return
	}

	player.Loses += increase

	_, err = s.tx.
		Model(&player).
		Column("loses").
		WherePK().
		Update()

	return
}

// PlayerIncreaseTies implements the `store.Transaction` interface.
func (s *Transaction) PlayerIncreaseTies(name string, increase uint) (err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}
	}()

	player := Player{ID: name}
	if err = s.tx.Select(&player); err != nil {
		if err == pg.ErrNoRows {
			return store.DontExistError(fmt.Sprintf(`The player %s doesn't exist`, name))
		}

		return
	}

	player.Ties += increase

	_, err = s.tx.
		Model(&player).
		Column("ties").
		WherePK().
		Update()

	return
}
