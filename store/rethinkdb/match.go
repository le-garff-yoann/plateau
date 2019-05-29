package rethinkdb

import (
	"fmt"
	"plateau/protocol"
	"plateau/store"
	"reflect"
	"sync"
	"time"

	rethinkdb "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

type match struct {
	ID string `rethinkdb:"id,omitempty"`

	CreatedAt time.Time  `rethinkdb:"created_at"`
	EndedAt   *time.Time `rethinkdb:"ended_at"`

	ConnectedPlayers []player `rethinkdb:"connected_player_ids,reference" rethinkdb_ref:"id"`

	NumberOfPlayersRequired uint     `rethinkdb:"number_of_players_required"`
	Players                 []player `rethinkdb:"player_ids,reference" rethinkdb_ref:"id"`

	Transactions []transaction `rethinkdb:"transactions"`
}

func matchFromProtocolStruct(m *protocol.Match) *match {
	var (
		connectedPlayers, players []player

		transactions []transaction
	)

	for _, p := range m.ConnectedPlayers {
		players = append(connectedPlayers, *playerFromProtocolStruct(&p))
	}

	for _, p := range m.Players {
		players = append(players, *playerFromProtocolStruct(&p))
	}

	for _, trx := range m.Transactions {
		transactions = append(transactions, *transactionFromProtocolStruct(&trx))
	}

	return &match{
		ID:                      m.ID,
		CreatedAt:               m.CreatedAt,
		EndedAt:                 m.EndedAt,
		ConnectedPlayers:        connectedPlayers,
		NumberOfPlayersRequired: m.NumberOfPlayersRequired,
		Players:                 players,
		Transactions:            transactions,
	}
}

func (s *match) toProtocolStruct() *protocol.Match {
	var (
		connectedPlayers, players []protocol.Player

		transactions []protocol.Transaction
	)

	for _, p := range s.ConnectedPlayers {
		players = append(connectedPlayers, *p.toProtocolStruct())
	}

	for _, p := range s.Players {
		players = append(players, *p.toProtocolStruct())
	}

	for _, trx := range s.Transactions {
		transactions = append(transactions, *trx.toProtocolStruct())
	}

	return &protocol.Match{
		ID:                      s.ID,
		CreatedAt:               s.CreatedAt,
		EndedAt:                 s.EndedAt,
		ConnectedPlayers:        connectedPlayers,
		NumberOfPlayersRequired: s.NumberOfPlayersRequired,
		Players:                 players,
		Transactions:            transactions,
	}
}

// matchStore ...
type matchStore struct {
	queryExecuter rethinkdb.QueryExecutor
}

func (s *matchStore) tableName() string {
	return "matchs"
}

func (s *matchStore) mergePredicateFunc() func(p rethinkdb.Term) interface{} {
	var playerStore playerStore

	return func(p rethinkdb.Term) interface{} {
		return map[string]interface{}{
			"connected_player_ids": rethinkdb.
				Table(playerStore.tableName()).
				GetAll(rethinkdb.Args(p.Field("connected_player_ids"))).
				CoerceTo("array"),
			"player_ids": rethinkdb.
				Table(playerStore.tableName()).
				GetAll(rethinkdb.Args(p.Field("player_ids"))).
				CoerceTo("array"),
			"transactions": p.
				Field("transactions").
				Map(func(p rethinkdb.Term) interface{} {
					return p.
						Merge(func(p rethinkdb.Term) interface{} {
							return map[string]interface{}{
								"holder_id": rethinkdb.
									Table(playerStore.tableName()).
									Get(p.Field("holder_id")),
							}
						})
				}),
		}
	}
}

func (s *matchStore) listTerm() rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		Map(func(p rethinkdb.Term) interface{} {
			return p.Field("id")
		})
}

func (s *matchStore) createTerm(m *protocol.Match) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		Insert(matchFromProtocolStruct(m))
}

func (s *matchStore) readTerm(id string) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(id).
		Merge(s.mergePredicateFunc())
}

func (s *matchStore) endedAtTerm(id string, val *time.Time) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(id).
		Update(map[string]interface{}{
			"ended_at": val,
		})
}

func (s *matchStore) connectPlayerTerm(id string, name string) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(id).
		Update(func(p rethinkdb.Term) interface{} {
			return map[string]interface{}{
				"connected_player_ids": rethinkdb.Branch(
					p.
						Field("connected_player_ids").
						Count(func(p rethinkdb.Term) interface{} {
							return p.Eq(name)
						}).
						Eq(0),
					p.
						Field("connected_player_ids").
						Append(name),
					p.
						Field("connected_player_ids"),
				),
			}
		})
}

func (s *matchStore) disconnectPlayerTerm(id string, name string) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(id).
		Update(func(p rethinkdb.Term) interface{} {
			return map[string]interface{}{
				"connected_player_ids": p.
					Field("connected_player_ids").
					Filter(func(p rethinkdb.Term) interface{} {
						return p.Ne(name)
					}),
			}
		})
}

func (s *matchStore) playerJoinsTerm(id string, name string) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(id).
		Update(func(p rethinkdb.Term) interface{} {
			return map[string]interface{}{
				"player_ids": rethinkdb.Branch(
					p.
						Field("player_ids").
						Count().
						Lt(
							p.
								Field("number_of_players_required"),
						).
						And(
							p.
								Field("player_ids").
								Count(func(p rethinkdb.Term) interface{} {
									return p.Eq(name)
								}).
								Eq(0),
						),
					p.
						Field("player_ids").
						Append(name),
					p.
						Field("player_ids"),
				),
			}
		})
}

func (s *matchStore) playerLeavesTerm(id string, name string) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(id).
		Update(func(p rethinkdb.Term) interface{} {
			return map[string]interface{}{
				"player_ids": rethinkdb.Branch(
					p.
						Field("player_ids").
						Count(func(p rethinkdb.Term) interface{} {
							return p.Eq(name)
						}).
						Gt(0),
					p.
						Field("player_ids").
						Filter(func(p rethinkdb.Term) interface{} {
							return p.Ne(name)
						}),
					p.
						Field("player_ids"),
				),
			}
		})
}

func (s *matchStore) createTransactionTerm(id string, trx *protocol.Transaction) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(id).
		Update(func(p rethinkdb.Term) interface{} {
			return map[string]interface{}{
				"transactions": p.
					Field("transactions").
					Append(transactionFromProtocolStruct(trx)),
			}
		})
}

func (s *matchStore) updateCurrentTransactionHolderTerm(id, holderName string) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(id).
		Update(func(p rethinkdb.Term) interface{} {
			return map[string]interface{}{
				"transactions": p.
					Field("transactions").
					ChangeAt(-1, p.
						Field("transactions").
						Nth(-1).
						Merge(map[string]interface{}{
							"holder_id": holderName,
						}),
					),
			}
		})
}

func (s *matchStore) addMessageToCurrentTransaction(id string, msg *protocol.Message) rethinkdb.Term {
	return rethinkdb.
		Table(s.tableName()).
		GetAll(id).
		Update(func(p rethinkdb.Term) interface{} {
			return map[string]interface{}{
				"transactions": p.
					Field("transactions").
					ChangeAt(-1, p.
						Field("transactions").
						Nth(-1).
						Merge(func(p rethinkdb.Term) interface{} {
							return map[string]interface{}{
								"messages": p.
									Field("messages").
									Append(messageFromProtocolStruct(msg)),
							}
						}),
					),
			}
		})
}

func (s *matchStore) matchChangesTerm(id string) rethinkdb.Term {
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

// List implements `store.matchStore` interface.
func (s *matchStore) List() ([]string, error) {
	cursor, err := s.listTerm().Run(s.queryExecuter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var IDs []string

	err = cursor.All(&IDs)
	return IDs, err
}

// Create  implements `store.matchStore` interface.
func (s *matchStore) Create(m protocol.Match) (string, error) {
	wRes, err := s.createTerm(&m).RunWrite(s.queryExecuter)
	if err != nil {
		return "", err
	}

	if len(wRes.GeneratedKeys) > 0 {
		return wRes.GeneratedKeys[0], nil
	}

	return m.ID, nil
}

// Read  implements `store.matchStore` interface.
func (s *matchStore) Read(id string) (*protocol.Match, error) {
	cursor, err := s.readTerm(id).Run(s.queryExecuter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var match match

	err = cursor.One(&match)
	if err != nil {
		if err == rethinkdb.ErrEmptyResult {
			return nil, store.DontExistError(err.Error())
		}

		return nil, err
	}

	return match.toProtocolStruct(), nil
}

// EndedAt implements `store.matchStore` interface.
func (s *matchStore) EndedAt(id string, val time.Time) error {
	_, err := s.endedAtTerm(id, &val).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}

	return err
}

// ConnectPlayer implements `store.matchStore` interface.
func (s *matchStore) ConnectPlayer(id, name string) error {
	wRes, err := s.connectPlayerTerm(id, name).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}
	if err != nil {
		return err
	}

	if wRes.Replaced == 0 {
		return store.PlayerConnectionError(fmt.Sprintf(`The player "%s" is already connected to the match "%s"`, name, id))
	}

	return nil
}

// DisconnectPlayer implements `store.matchStore` interface.
func (s *matchStore) DisconnectPlayer(id, name string) error {
	_, err := s.disconnectPlayerTerm(id, name).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}

	return err
}

// PlayerJoins implements `store.matchStore` interface.
func (s *matchStore) PlayerJoins(id, name string) error {
	wRes, err := s.playerJoinsTerm(id, name).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}
	if err != nil {
		return err
	}

	if wRes.Replaced == 0 {
		return store.PlayerParticipationError(fmt.Sprintf(`There are no more spot in match %s | The player "%s" is already in`, id, name))
	}

	return nil
}

// PlayerLeaves implements `store.matchStore` interface.
func (s *matchStore) PlayerLeaves(id, name string) error {
	wRes, err := s.playerLeavesTerm(id, name).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}
	if err != nil {
		return err
	}

	if wRes.Replaced == 0 {
		return store.PlayerParticipationError(fmt.Sprintf(`The player "%s" is already out of the match %s`, name, id))
	}

	return nil
}

// CreateTransaction implements `store.matchStore` interface.
func (s *matchStore) CreateTransaction(id string, trx protocol.Transaction) error {
	_, err := s.createTransactionTerm(id, &trx).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}

	return err
}

// UpdateCurrentTransactionHolder implements `store.matchStore` interface.
func (s *matchStore) UpdateCurrentTransactionHolder(id, newHolderName string) error {
	_, err := s.updateCurrentTransactionHolderTerm(id, newHolderName).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}

	return err
}

// AddMessageToCurrentTransaction implements `store.matchStore` interface.
func (s *matchStore) AddMessageToCurrentTransaction(id string, msg protocol.Message) error {
	_, err := s.addMessageToCurrentTransaction(id, &msg).RunWrite(s.queryExecuter)
	if err == rethinkdb.ErrEmptyResult {
		return store.DontExistError(err.Error())
	}

	return err
}

// CreateTransactionsChangeIterator implements `store.matchStore` interface.
func (s *matchStore) CreateTransactionsChangeIterator(id string) (store.TransactionChangeIterator, error) {
	cursor, err := s.matchChangesTerm(id).Run(s.queryExecuter)
	if err != nil {
		if err == rethinkdb.ErrEmptyResult {
			return nil, store.DontExistError(err.Error())
		}

		return nil, err
	}

	return &TransactionChangeIterator{cursor: cursor, trxChanges: []store.TransactionChange{}}, nil
}

type matchChangeResponse struct {
	NewValue match `rethinkdb:"new_val"`
	OldValue match `rethinkdb:"old_val"`
}

// TransactionChangeIterator implements `store.TransactionChangeIterator` interface.
type TransactionChangeIterator struct {
	mux sync.Mutex

	cursor     *rethinkdb.Cursor
	trxChanges []store.TransactionChange
}

// Next implements `store.TransactionChangeIterator` interface.
func (s *TransactionChangeIterator) Next(trxChange *store.TransactionChange) bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	for len(s.trxChanges) == 0 {
		var changeRes matchChangeResponse
		if b := s.cursor.Next(&changeRes); !b {
			return b
		}

		for i := range changeRes.NewValue.Transactions {
			if i >= len(changeRes.OldValue.Transactions) {
				s.trxChanges = append(s.trxChanges, store.TransactionChange{
					Old: nil,
					New: changeRes.NewValue.Transactions[i].toProtocolStruct(),
				})
			} else {
				var (
					oldTrx = changeRes.OldValue.Transactions[i].toProtocolStruct()
					newTrx = changeRes.NewValue.Transactions[i].toProtocolStruct()
				)

				if !reflect.DeepEqual(oldTrx, newTrx) {
					s.trxChanges = append(s.trxChanges, store.TransactionChange{Old: oldTrx, New: newTrx})
				}
			}
		}
	}

	*trxChange = s.trxChanges[0]
	s.trxChanges = s.trxChanges[1:]

	return true
}

// Close implements `store.TransactionChangeIterator` interface.
func (s *TransactionChangeIterator) Close() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.cursor.Close()
}
