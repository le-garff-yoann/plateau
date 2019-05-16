package rethinkdb

import (
	"plateau/store"

	rethinkdb "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

// Player ...
type Player struct {
	Name string `rethinkdb:"id"`

	Password string `rethinkdb:"password"`

	Wins  uint `rethinkdb:"wins"`
	Loses uint `rethinkdb:"loses"`
	Ties  uint `rethinkdb:"ties"`
}

func playerFromStoreStruct(p store.Player) *Player {
	return &Player{
		p.Name, p.Password,
		p.Wins, p.Loses, p.Ties,
	}
}

func (s *Player) toStoreStruct() *store.Player {
	return &store.Player{
		Name:     s.Name,
		Password: s.Password,
		Wins:     s.Wins,
		Loses:    s.Loses,
		Ties:     s.Ties,
	}
}

// PlayerStore ...
type PlayerStore struct {
	queryExecuter rethinkdb.QueryExecutor
}

func (s *PlayerStore) tableName() string {
	return "players"
}

func (s *PlayerStore) listTerm() rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		Map(func(p rethinkdb.Term) interface{} {
			return p.Field("id")
		})
}

func (s *PlayerStore) createTerm(p store.Player) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		Insert(playerFromStoreStruct(p))
}

func (s *PlayerStore) readTerm(name string) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(name)
}

func (s *PlayerStore) increaseScoreTerm(field string, name string, val uint) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(name).
		Update(func(p rethinkdb.Term) interface{} {
			return map[string]interface{}{
				field: p.
					Field(field).
					Add(val),
			}
		})
}

// List ...
func (s *PlayerStore) List() ([]string, error) {
	cursor, err := s.listTerm().Run(s.queryExecuter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var names []string

	err = cursor.All(&names)
	return names, err
}

// Create ...
func (s *PlayerStore) Create(p store.Player) error {
	_, err := s.createTerm(p).RunWrite(s.queryExecuter)
	if rethinkdb.IsConflictErr(err) {
		return store.DuplicateError(err.Error())
	}

	return err
}

// Read ...
func (s *PlayerStore) Read(name string) (*store.Player, error) {
	cursor, err := s.readTerm(name).Run(s.queryExecuter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var player Player

	err = cursor.One(&player)
	if err != nil {
		if err == rethinkdb.ErrEmptyResult {
			return nil, store.DontExistError(err.Error())
		}

		return nil, err
	}

	return player.toStoreStruct(), nil
}

// IncreaseWins ...
func (s *PlayerStore) IncreaseWins(name string, val uint) error {
	_, err := s.increaseScoreTerm("wins", name, val).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}

	return err
}

// IncreaseLoses ...
func (s *PlayerStore) IncreaseLoses(name string, val uint) error {
	_, err := s.increaseScoreTerm("loses", name, val).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}

	return err
}

// IncreaseTies ...
func (s *PlayerStore) IncreaseTies(name string, val uint) error {
	_, err := s.increaseScoreTerm("ties", name, val).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}

	return err
}
