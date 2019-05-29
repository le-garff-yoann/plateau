package rethinkdb

import (
	"plateau/protocol"
	"plateau/store"

	rethinkdb "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

type player struct {
	Name     string `rethinkdb:"id"`
	Password string `rethinkdb:"password"`

	Wins  uint `rethinkdb:"wins"`
	Loses uint `rethinkdb:"loses"`
	Ties  uint `rethinkdb:"ties"`
}

func playerFromProtocolStruct(p *protocol.Player) *player {
	return &player{
		p.Name, p.Password,
		p.Wins, p.Loses, p.Ties,
	}
}

func (s *player) toProtocolStruct() *protocol.Player {
	return &protocol.Player{
		Name:     s.Name,
		Password: s.Password,
		Wins:     s.Wins,
		Loses:    s.Loses,
		Ties:     s.Ties,
	}
}

// playerStore ...
type playerStore struct {
	queryExecuter rethinkdb.QueryExecutor
}

func (s *playerStore) tableName() string {
	return "players"
}

func (s *playerStore) listTerm() rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		Map(func(p rethinkdb.Term) interface{} {
			return p.Field("id")
		})
}

func (s *playerStore) createTerm(p *protocol.Player) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		Insert(playerFromProtocolStruct(p))
}

func (s *playerStore) readTerm(name string) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(name)
}

func (s *playerStore) connectedTerm(name string, val bool) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(name).
		Update(map[string]interface{}{
			"connected": val,
		})
}

func (s *playerStore) increaseScoreTerm(field, name string, increase uint) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(name).
		Update(func(p rethinkdb.Term) interface{} {
			return map[string]interface{}{
				field: p.
					Field(field).
					Add(increase),
			}
		})
}

// List implements `store.playerStore` interface.
func (s *playerStore) List() ([]string, error) {
	cursor, err := s.listTerm().Run(s.queryExecuter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var names []string

	err = cursor.All(&names)
	return names, err
}

// Create implements `store.playerStore` interface.
func (s *playerStore) Create(p protocol.Player) error {
	_, err := s.createTerm(&p).RunWrite(s.queryExecuter)
	if rethinkdb.IsConflictErr(err) {
		return store.DuplicateError(err.Error())
	}

	return err
}

// Read implements `store.playerStore` interface.
func (s *playerStore) Read(name string) (*protocol.Player, error) {
	cursor, err := s.readTerm(name).Run(s.queryExecuter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var player player

	err = cursor.One(&player)
	if err != nil {
		if err == rethinkdb.ErrEmptyResult {
			return nil, store.DontExistError(err.Error())
		}

		return nil, err
	}

	return player.toProtocolStruct(), nil
}

// IncreaseWins implements `store.playerStore` interface.
func (s *playerStore) IncreaseWins(name string, increase uint) error {
	_, err := s.increaseScoreTerm("wins", name, increase).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}

	return err
}

// IncreaseLoses implements `store.playerStore` interface.
func (s *playerStore) IncreaseLoses(name string, increase uint) error {
	_, err := s.increaseScoreTerm("loses", name, increase).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}

	return err
}

// IncreaseTies implements `store.playerStore` interface.
func (s *playerStore) IncreaseTies(name string, increase uint) error {
	_, err := s.increaseScoreTerm("ties", name, increase).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}

	return err
}
