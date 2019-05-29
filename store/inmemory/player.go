package inmemory

import (
	"fmt"
	"plateau/protocol"
	"plateau/store"

	"github.com/ulule/deepcopier"
)

type playerStore struct {
	*inMemory
}

// List ...
func (s *playerStore) List() (names []string, err error) {
	s.inMemory.mux.RLock()
	defer s.inMemory.mux.RUnlock()

	for _, p := range s.inMemory.players {
		names = append(names, p.Name)
	}

	return names, nil
}

// Create ...
func (s *playerStore) Create(p protocol.Player) error {
	s.inMemory.mux.Lock()
	defer s.inMemory.mux.Unlock()

	for _, player := range s.inMemory.players {
		if p.Name == player.Name {
			return store.DuplicateError(fmt.Sprintf(`The player %s already exist`, p.Name))
		}
	}

	s.inMemory.players = append(s.inMemory.players, &p)

	return nil
}

// Read ...
func (s *playerStore) Read(name string) (*protocol.Player, error) {
	s.inMemory.mux.RLock()
	defer s.inMemory.mux.RUnlock()

	p := s.inMemory.player(name)
	if p == nil {
		return nil, store.DontExistError(fmt.Sprintf(`The player %s doesn't exist`, name))
	}

	var copy protocol.Player
	deepcopier.Copy(p).To(&copy)

	return &copy, nil
}

// IncreaseWins ...
func (s *playerStore) IncreaseWins(name string, increase uint) error {
	s.inMemory.mux.Lock()
	defer s.inMemory.mux.Unlock()

	p := s.inMemory.player(name)
	if p == nil {
		return store.DontExistError(fmt.Sprintf(`The player %s doesn't exist`, name))
	}

	p.Wins = p.Wins + increase

	return nil
}

// IncreaseLoses ...
func (s *playerStore) IncreaseLoses(name string, increase uint) error {
	s.inMemory.mux.Lock()
	defer s.inMemory.mux.Unlock()

	p := s.inMemory.player(name)
	if p == nil {
		return store.DontExistError(fmt.Sprintf(`The player %s doesn't exist`, name))
	}

	p.Loses = p.Loses + increase

	return nil
}

// IncreaseTies ...
func (s *playerStore) IncreaseTies(name string, increase uint) error {
	s.inMemory.mux.Lock()
	defer s.inMemory.mux.Unlock()

	p := s.inMemory.player(name)
	if p == nil {
		return store.DontExistError(fmt.Sprintf(`The player %s doesn't exist`, name))
	}

	p.Ties = p.Ties + increase

	return nil
}
