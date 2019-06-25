package inmemory

import (
	"fmt"
	"plateau/protocol"
	"plateau/store"

	"github.com/ulule/deepcopier"
)

// PlayerList ...
func (s *Transaction) PlayerList() (names []string, err error) {
	defer func() {
		s.errors = append(s.errors, err)
	}()

	for _, p := range s.inMemoryCopy.Players {
		names = append(names, p.Name)
	}

	return names, nil
}

// PlayerCreate ...
func (s *Transaction) PlayerCreate(p protocol.Player) (err error) {
	defer func() {
		s.errors = append(s.errors, err)
	}()

	for _, player := range s.inMemoryCopy.Players {
		if p.Name == player.Name {
			return store.DuplicateError(fmt.Sprintf(`The player %s already exist`, p.Name))
		}
	}

	s.inMemoryCopy.Players = append(s.inMemoryCopy.Players, &p)

	return nil
}

// PlayerRead ...
func (s *Transaction) PlayerRead(name string) (_ *protocol.Player, err error) {
	defer func() {
		s.errors = append(s.errors, err)
	}()

	p := s.inMemoryCopy.player(name)
	if p == nil {
		return nil, store.DontExistError(fmt.Sprintf(`The player %s doesn't exist`, name))
	}

	var playerCopy protocol.Player
	deepcopier.Copy(p).To(&playerCopy)

	return &playerCopy, nil
}

// PlayerIncreaseWins ...
func (s *Transaction) PlayerIncreaseWins(name string, increase uint) (err error) {
	defer func() {
		s.errors = append(s.errors, err)
	}()

	p := s.inMemoryCopy.player(name)
	if p == nil {
		return store.DontExistError(fmt.Sprintf(`The player %s doesn't exist`, name))
	}

	p.Wins = p.Wins + increase

	return nil
}

// PlayerIncreaseLoses ...
func (s *Transaction) PlayerIncreaseLoses(name string, increase uint) (err error) {
	defer func() {
		s.errors = append(s.errors, err)
	}()

	p := s.inMemoryCopy.player(name)
	if p == nil {
		return store.DontExistError(fmt.Sprintf(`The player %s doesn't exist`, name))
	}

	p.Loses = p.Loses + increase

	return nil
}

// PlayerIncreaseTies ...
func (s *Transaction) PlayerIncreaseTies(name string, increase uint) (err error) {
	defer func() {
		s.errors = append(s.errors, err)
	}()

	p := s.inMemoryCopy.player(name)
	if p == nil {
		return store.DontExistError(fmt.Sprintf(`The player %s doesn't exist`, name))
	}

	p.Ties = p.Ties + increase

	return nil
}
