package rethinkdb

import (
	"plateau/store"
	"time"

	rethinkdb "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

// Match ...
type Match struct {
	ID string `rethinkdb:"id,omitempty"`

	CreatedAt time.Time `rethinkdb:"created_at"`
	EndedAt   time.Time `rethinkdb:"ended_at"`

	NumberOfPlayersRequired uint     `rethinkdb:"number_of_players_required"`
	Players                 []Player `rethinkdb:"player_ids,reference" rethinkdb_ref:"id"`

	Running bool `rethinkdb:"running"`

	EventContainers []EventContainer `rethinkdb:"event_containers"`
}

func matchFromStoreStruct(m store.Match) *Match {
	var (
		players []Player

		createdAt, endedAt time.Time

		eventContainers []EventContainer
	)

	if m.CreatedAt != nil {
		createdAt = *m.CreatedAt
	}

	if m.EndedAt != nil {
		endedAt = *m.EndedAt
	}

	for _, p := range m.Players {
		players = append(players, *playerFromStoreStruct(*p))
	}

	for _, ec := range m.EventContainers {
		eventContainers = append(eventContainers, *eventContainerFromStoreStruct(*ec))
	}

	return &Match{
		ID:                      m.ID,
		CreatedAt:               createdAt,
		EndedAt:                 endedAt,
		NumberOfPlayersRequired: m.NumberOfPlayersRequired,
		Players:                 players,
		Running:                 m.Running,
		EventContainers:         eventContainers,
	}
}

func (s *Match) toStoreStruct() *store.Match {
	var (
		players []*store.Player

		eventContainers []*store.EventContainer
	)

	for _, p := range s.Players {
		players = append(players, p.toStoreStruct())
	}

	for _, ec := range s.EventContainers {
		eventContainers = append(eventContainers, ec.toStoreStruct())
	}

	return &store.Match{
		ID:                      s.ID,
		CreatedAt:               &s.CreatedAt,
		EndedAt:                 &s.EndedAt,
		NumberOfPlayersRequired: s.NumberOfPlayersRequired,
		Players:                 players,
		Running:                 s.Running,
		EventContainers:         eventContainers,
	}
}

// MatchStore ...
type MatchStore struct {
	queryExecuter rethinkdb.QueryExecutor
}

func (s *MatchStore) tableName() string {
	return "matchs"
}

type matchChangeResponse struct {
	NewValue Match `rethinkdb:"new_val"`
	OldValue Match `rethinkdb:"old_val"`
}

func (s *MatchStore) mergePredicateFunc() func(p rethinkdb.Term) interface{} {
	var playerStore PlayerStore

	return func(p rethinkdb.Term) interface{} {
		return map[string]interface{}{
			"player_ids": rethinkdb.
				Table(playerStore.tableName()).
				GetAll(rethinkdb.Args(p.Field("player_ids"))).
				CoerceTo("array"),
			"event_containers": p.
				Field("event_containers").
				Map(func(pp rethinkdb.Term) interface{} {
					return pp.
						Merge(func(ppp rethinkdb.Term) interface{} {
							return map[string]interface{}{
								"emitter_id": rethinkdb.
									Table(playerStore.tableName()).
									Get(ppp.Field("emitter_id")),
								"receiver_ids": rethinkdb.
									Table(playerStore.tableName()).
									GetAll(rethinkdb.Args(ppp.Field("receiver_ids"))).
									CoerceTo("array"),
								"subject_ids": rethinkdb.
									Table(playerStore.tableName()).
									GetAll(rethinkdb.Args(ppp.Field("subject_ids"))).
									CoerceTo("array"),
							}
						})
				}).CoerceTo("array"),
		}
	}
}

func (s *MatchStore) listTerm() rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		Map(func(p rethinkdb.Term) interface{} {
			return p.Field("id")
		})
}

func (s *MatchStore) createTerm(m store.Match) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		Insert(matchFromStoreStruct(m))
}

func (s *MatchStore) readTerm(id string) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(id).
		Merge(s.mergePredicateFunc())
}

func (s *MatchStore) endedAtTerm(id string, val *time.Time) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(id).
		Update(map[string]interface{}{
			"ended_at": *val,
		})
}

func (s *MatchStore) addPlayerTerm(id string, name string) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(id).
		Update(func(p rethinkdb.Term) interface{} {
			return map[string]interface{}{
				"player_ids": p.
					Field("player_ids").
					Append(name).
					Distinct(),
			}
		})
}

func (s *MatchStore) removePlayerTerm(id string, name string) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(id).
		Update(func(p rethinkdb.Term) interface{} {
			return map[string]interface{}{
				"player_ids": p.
					Field("player_ids").
					Filter(func(pp rethinkdb.Term) interface{} {
						return pp.Ne(name)
					}),
			}
		})
}

func (s *MatchStore) runningTerm(id string, val bool) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(id).
		Update(map[string]interface{}{
			"running": val,
		})
}

func (s *MatchStore) createEventContainerTerm(id string, ec store.EventContainer) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(id).
		Update(func(p rethinkdb.Term) interface{} {
			return map[string]interface{}{
				"event_containers": p.
					Field("event_containers").
					Append(*eventContainerFromStoreStruct(ec)),
			}
		})
}

func (s *MatchStore) eventContainerChangesTerm(id string) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(id).
		Changes().
		Merge(func(p rethinkdb.Term) interface{} {
			return map[string]interface{}{
				"new_val": p.
					Field("new_val").
					Merge(s.mergePredicateFunc()),
				"old_val": p.
					Field("old_val").
					Merge(s.mergePredicateFunc()),
			}
		})
}

// List ...
func (s *MatchStore) List() ([]string, error) {
	cursor, err := s.listTerm().Run(s.queryExecuter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var IDs []string

	err = cursor.All(&IDs)
	return IDs, err
}

// Create ...
func (s *MatchStore) Create(m store.Match) (string, error) {
	wRes, err := s.createTerm(m).RunWrite(s.queryExecuter)
	if err != nil {
		return "", err
	}

	if len(wRes.GeneratedKeys) > 0 {
		return wRes.GeneratedKeys[0], nil
	}

	return "", nil
}

// Read ...
func (s *MatchStore) Read(id string) (*store.Match, error) {
	cursor, err := s.readTerm(id).Run(s.queryExecuter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var match Match

	err = cursor.One(&match)
	if err != nil {
		if err == rethinkdb.ErrEmptyResult {
			return nil, store.DontExistError(err.Error())
		}

		return nil, err
	}

	return match.toStoreStruct(), nil
}

// EndedAt ...
func (s *MatchStore) EndedAt(id string, val *time.Time) error {
	_, err := s.endedAtTerm(id, val).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}

	return err
}

// AddPlayer ...
func (s *MatchStore) AddPlayer(id string, name string) error {
	_, err := s.addPlayerTerm(id, name).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}

	return err
}

// RemovePlayer ...
func (s *MatchStore) RemovePlayer(id string, name string) error {
	_, err := s.removePlayerTerm(id, name).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}

	return err
}

// Running ...
func (s *MatchStore) Running(id string, val bool) error {
	_, err := s.runningTerm(id, val).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}

	return err
}

// CreateEventContainer ...
func (s *MatchStore) CreateEventContainer(id string, ec store.EventContainer) error {
	_, err := s.createEventContainerTerm(id, ec).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}

	return err
}

// CreateEventContainerBroadcaster ...
func (s *MatchStore) CreateEventContainerBroadcaster(id string) (*store.EventContainerBroadcaster, error) {
	cursor, err := s.eventContainerChangesTerm(id).Run(s.queryExecuter)
	if err != nil {
		if err == rethinkdb.ErrEmptyResult {
			return nil, store.DontExistError(err.Error())
		}

		return nil, err
	}

	br := store.NewEventContainerBroadcaster()

	go func() {
		var changeRes matchChangeResponse

		for cursor.Next(&changeRes) {
			for i := len(changeRes.OldValue.EventContainers); i < len(changeRes.NewValue.EventContainers); i++ {
				br.Emitter <- *changeRes.NewValue.EventContainers[i].toStoreStruct()
			}
		}
	}()

	go func() {
		for {
			<-br.Done

			if cursor.Close() == nil {
				return
			}
		}
	}()

	return br, nil
}
