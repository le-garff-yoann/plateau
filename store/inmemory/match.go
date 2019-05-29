package inmemory

import (
	"fmt"
	"plateau/broadcaster"
	"plateau/protocol"
	"plateau/store"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/ulule/deepcopier"
)

type match struct {
	ID string

	CreatedAt time.Time
	EndedAt   *time.Time

	ConnectedPlayers map[string]protocol.Player

	NumberOfPlayersRequired uint
	Players                 map[string]protocol.Player

	Transactions []transaction
}

func matchFromProtocolStruct(m *protocol.Match) *match {
	var (
		connectedplayers = make(map[string]protocol.Player)
		players          = make(map[string]protocol.Player)

		transactions []transaction
	)

	for _, p := range m.ConnectedPlayers {
		connectedplayers[p.Name] = p
	}

	for _, p := range m.Players {
		players[p.Name] = p
	}

	for _, trx := range m.Transactions {
		transactions = append(transactions, *transactionFromProtocolStruct(&trx))
	}

	return &match{
		ID:                      m.ID,
		CreatedAt:               m.CreatedAt,
		EndedAt:                 m.EndedAt,
		ConnectedPlayers:        connectedplayers,
		NumberOfPlayersRequired: m.NumberOfPlayersRequired,
		Players:                 players,
		Transactions:            transactions,
	}
}

func (s *match) toProtocolStruct(pPlayers []*protocol.Player) *protocol.Match {
	var (
		connectedplayers, players []protocol.Player

		pPlayersMap = make(map[string]protocol.Player)

		transactions []protocol.Transaction
	)

	for _, p := range pPlayers {
		pPlayersMap[p.Name] = *p
	}

	for pName, p := range s.ConnectedPlayers {
		pp, ok := pPlayersMap[pName]
		if ok {
			connectedplayers = append(connectedplayers, pp)
		} else {
			connectedplayers = append(connectedplayers, p)
		}
	}

	for pName, p := range s.Players {
		pp, ok := pPlayersMap[pName]
		if ok {
			players = append(players, pp)
		} else {
			players = append(players, p)
		}
	}

	for _, trx := range s.Transactions {
		transactions = append(transactions, *trx.toProtocolStruct(pPlayers))
	}

	return &protocol.Match{
		ID:                      s.ID,
		CreatedAt:               s.CreatedAt,
		EndedAt:                 s.EndedAt,
		ConnectedPlayers:        connectedplayers,
		NumberOfPlayersRequired: s.NumberOfPlayersRequired,
		Players:                 players,
		Transactions:            transactions,
	}
}

type matchStore struct {
	*inMemory

	trxChangesBroadcaster *broadcaster.Broadcaster
}

func newmatchStore(inm *inMemory) *matchStore {
	br := broadcaster.New()

	go br.Run()

	return &matchStore{inm, br}
}

func (s *matchStore) close() error {
	s.trxChangesBroadcaster.Done()

	return nil
}

// List ...
func (s *matchStore) List() (IDs []string, err error) {
	s.inMemory.mux.RLock()
	defer s.inMemory.mux.RUnlock()

	for _, m := range s.inMemory.matchs {
		IDs = append(IDs, m.ID)
	}

	return IDs, nil
}

// Create ...
func (s *matchStore) Create(m protocol.Match) (id string, err error) {
	s.inMemory.mux.Lock()
	defer s.inMemory.mux.Unlock()

	m.ID = uuid.NewV4().String()

	s.inMemory.matchs = append(s.inMemory.matchs, matchFromProtocolStruct(&m))

	return m.ID, nil
}

// Read ...
func (s *matchStore) Read(id string) (*protocol.Match, error) {
	s.inMemory.mux.RLock()
	defer s.inMemory.mux.RUnlock()

	m := s.inMemory.match(id)
	if m == nil {
		return nil, store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	var copy protocol.Match
	deepcopier.Copy(m.toProtocolStruct(s.players)).To(&copy)

	return &copy, nil
}

// EndedAt ...
func (s *matchStore) EndedAt(id string, val time.Time) error {
	s.inMemory.mux.Lock()
	defer s.inMemory.mux.Unlock()

	m := s.inMemory.match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	m.EndedAt = &val

	return nil
}

// CreateTransaction ...
func (s *matchStore) CreateTransaction(id string, trx protocol.Transaction) error {
	s.inMemory.mux.Lock()
	defer s.inMemory.mux.Unlock()

	m := s.inMemory.match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	m.Transactions = append(m.Transactions, *transactionFromProtocolStruct(&trx))

	s.trxChangesBroadcaster.Submit(store.TransactionChange{Old: nil, New: &trx})

	return nil
}

// UpdateCurrentTransactionHolder ...
func (s *matchStore) UpdateCurrentTransactionHolder(id, newHolderName string) error {
	s.inMemory.mux.Lock()
	defer s.inMemory.mux.Unlock()

	m := s.inMemory.match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	var oldTrx transaction
	deepcopier.Copy(m.Transactions[len(m.Transactions)-1]).To(&oldTrx)

	m.Transactions[len(m.Transactions)-1].Holder.Name = newHolderName

	s.trxChangesBroadcaster.Submit(store.TransactionChange{
		Old: oldTrx.toProtocolStruct(s.players),
		New: m.Transactions[len(m.Transactions)-1].toProtocolStruct(s.players),
	})

	return nil
}

// AddMessageToCurrentTransaction ...
func (s *matchStore) AddMessageToCurrentTransaction(id string, message protocol.Message) error {
	s.inMemory.mux.Lock()
	defer s.inMemory.mux.Unlock()

	m := s.inMemory.match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	var oldTrx transaction
	deepcopier.Copy(m.Transactions[len(m.Transactions)-1]).To(&oldTrx)

	trx := &m.Transactions[len(m.Transactions)-1]
	trx.Messages = append(trx.Messages, message)

	s.trxChangesBroadcaster.Submit(store.TransactionChange{
		Old: oldTrx.toProtocolStruct(s.players),
		New: trx.toProtocolStruct(s.players),
	})

	return nil
}

// ConnectPlayer ...
func (s *matchStore) ConnectPlayer(id, name string) error {
	s.inMemory.mux.Lock()
	defer s.inMemory.mux.Unlock()

	m := s.inMemory.match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	_, ok := m.ConnectedPlayers[name]
	if ok {
		return store.PlayerConnectionError(fmt.Sprintf(`The player "%s" is already connected to the match "%s"`, name, id))
	}

	m.ConnectedPlayers[name] = protocol.Player{Name: name}

	return nil
}

// DisconnectPlayer ...
func (s *matchStore) DisconnectPlayer(id, name string) error {
	s.inMemory.mux.Lock()
	defer s.inMemory.mux.Unlock()

	m := s.inMemory.match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	delete(m.ConnectedPlayers, name)

	return nil
}

// PlayerJoins ...
func (s *matchStore) PlayerJoins(id, name string) error {
	s.inMemory.mux.Lock()
	defer s.inMemory.mux.Unlock()

	m := s.inMemory.match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	if len(m.Players) >= int(m.NumberOfPlayersRequired) {
		return store.PlayerParticipationError(fmt.Sprintf(`There are no more spot in match %s`, id))
	}

	_, ok := m.Players[name]
	if ok {
		return store.PlayerParticipationError(fmt.Sprintf(`The player "%s" is already in the match %s`, name, id))
	}

	m.Players[name] = protocol.Player{Name: name}

	return nil
}

// PlayerLeaves ...
func (s *matchStore) PlayerLeaves(id, name string) error {
	s.inMemory.mux.Lock()
	defer s.inMemory.mux.Unlock()

	m := s.inMemory.match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	_, ok := m.Players[name]
	if !ok {
		return store.PlayerParticipationError(fmt.Sprintf(`The player "%s" is already out of the match %s`, name, id))
	}

	delete(m.Players, name)

	return nil
}

// CreateTransactionsChangeIterator ...
func (s *matchStore) CreateTransactionsChangeIterator(id string) (store.TransactionChangeIterator, error) {
	itr := TransactionChangeIterator{trxChangesBroadcaster: s.trxChangesBroadcaster}

	itr.trxChangesBroadcasterChan, itr.trxChangesBroadcasterUUID = s.trxChangesBroadcaster.Subscribe()

	return &itr, nil
}

// TransactionChangeIterator implements `store.TransactionChangeIterator` interface.
type TransactionChangeIterator struct {
	trxChangesBroadcaster *broadcaster.Broadcaster

	trxChangesBroadcasterChan <-chan interface{}
	trxChangesBroadcasterUUID uuid.UUID
}

// Next implements `store.TransactionChangeIterator` interface.
func (s *TransactionChangeIterator) Next(trxChange *store.TransactionChange) bool {
	v, ok := <-s.trxChangesBroadcasterChan
	if !ok {
		return false
	}

	*trxChange = v.(store.TransactionChange)

	return true
}

// Close implements `store.TransactionChangeIterator` interface.
func (s *TransactionChangeIterator) Close() error {
	s.trxChangesBroadcaster.Unsubscribe(s.trxChangesBroadcasterUUID)

	return nil
}
