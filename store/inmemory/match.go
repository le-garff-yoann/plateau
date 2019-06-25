package inmemory

import (
	"fmt"
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

	Deals []deal
}

func matchFromProtocolStruct(m *protocol.Match) *match {
	var (
		connectedplayers = make(map[string]protocol.Player)
		players          = make(map[string]protocol.Player)

		deals []deal
	)

	for _, p := range m.ConnectedPlayers {
		connectedplayers[p.Name] = p
	}

	for _, p := range m.Players {
		players[p.Name] = p
	}

	for _, deal := range m.Deals {
		deals = append(deals, *dealFromProtocolStruct(&deal))
	}

	return &match{
		ID:                      m.ID,
		CreatedAt:               m.CreatedAt,
		EndedAt:                 m.EndedAt,
		ConnectedPlayers:        connectedplayers,
		NumberOfPlayersRequired: m.NumberOfPlayersRequired,
		Players:                 players,
		Deals:                   deals,
	}
}

func (s *match) toProtocolStruct(pPlayers []*protocol.Player) *protocol.Match {
	var (
		connectedplayers, players []protocol.Player

		pPlayersMap = make(map[string]protocol.Player)

		deals []protocol.Deal
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

	for _, deal := range s.Deals {
		deals = append(deals, *deal.toProtocolStruct(pPlayers))
	}

	return &protocol.Match{
		ID:                      s.ID,
		CreatedAt:               s.CreatedAt,
		EndedAt:                 s.EndedAt,
		ConnectedPlayers:        connectedplayers,
		NumberOfPlayersRequired: s.NumberOfPlayersRequired,
		Players:                 players,
		Deals:                   deals,
	}
}

// MatchList ...
func (s *Transaction) MatchList() (IDs []string, err error) {
	defer func() {
		s.errors = append(s.errors, err)
	}()

	for _, m := range s.inMemoryCopy.Matchs {
		IDs = append(IDs, m.ID)
	}

	return IDs, nil
}

// MatchCreate ...
func (s *Transaction) MatchCreate(m protocol.Match) (id string, err error) {
	defer func() {
		s.errors = append(s.errors, err)
	}()

	m.ID = uuid.NewV4().String()

	s.inMemoryCopy.Matchs = append(s.inMemoryCopy.Matchs, matchFromProtocolStruct(&m))

	return m.ID, nil
}

// MatchRead ...
func (s *Transaction) MatchRead(id string) (_ *protocol.Match, err error) {
	defer func() {
		s.errors = append(s.errors, err)
	}()

	m := s.inMemoryCopy.match(id)
	if m == nil {
		return nil, store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	var matchCopy protocol.Match
	deepcopier.Copy(m.toProtocolStruct(s.inMemoryCopy.Players)).To(&matchCopy)

	return &matchCopy, nil
}

// MatchEndedAt ...
func (s *Transaction) MatchEndedAt(id string, val time.Time) (err error) {
	defer func() {
		s.errors = append(s.errors, err)
	}()

	m := s.inMemoryCopy.match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	m.EndedAt = &val

	return nil
}

// MatchCreateDeal ...
func (s *Transaction) MatchCreateDeal(id string, deal protocol.Deal) (err error) {
	defer func() {
		s.errors = append(s.errors, err)
	}()

	m := s.inMemoryCopy.match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	m.Deals = append(m.Deals, *dealFromProtocolStruct(&deal))

	s.dealChangeSubmitter(&store.DealChange{Old: nil, New: &deal})

	return nil
}

// MatchUpdateCurrentDealHolder ...
func (s *Transaction) MatchUpdateCurrentDealHolder(id, newHolderName string) (err error) {
	defer func() {
		s.errors = append(s.errors, err)
	}()

	m := s.inMemoryCopy.match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	var oldDeal deal
	deepcopier.Copy(m.Deals[len(m.Deals)-1]).To(&oldDeal)

	m.Deals[len(m.Deals)-1].Holder.Name = newHolderName

	s.dealChangeSubmitter(&store.DealChange{
		Old: oldDeal.toProtocolStruct(s.inMemoryCopy.Players),
		New: m.Deals[len(m.Deals)-1].toProtocolStruct(s.inMemoryCopy.Players),
	})

	return nil
}

// MatchAddMessageToCurrentDeal ...
func (s *Transaction) MatchAddMessageToCurrentDeal(id string, message protocol.Message) (err error) {
	defer func() {
		s.errors = append(s.errors, err)
	}()

	m := s.inMemoryCopy.match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	var oldDeal deal
	deepcopier.Copy(m.Deals[len(m.Deals)-1]).To(&oldDeal)

	deal := &m.Deals[len(m.Deals)-1]
	deal.Messages = append(deal.Messages, message)

	s.dealChangeSubmitter(&store.DealChange{
		Old: oldDeal.toProtocolStruct(s.inMemoryCopy.Players),
		New: deal.toProtocolStruct(s.inMemoryCopy.Players),
	}) // TODO-1

	return nil
}

// MatchConnectPlayer ...
func (s *Transaction) MatchConnectPlayer(id, name string) (err error) {
	defer func() {
		s.errors = append(s.errors, err)
	}()

	m := s.inMemoryCopy.match(id)
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

// MatchDisconnectPlayer ...
func (s *Transaction) MatchDisconnectPlayer(id, name string) (err error) {
	defer func() {
		s.errors = append(s.errors, err)
	}()

	m := s.inMemoryCopy.match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	delete(m.ConnectedPlayers, name)

	return nil
}

// MatchPlayerJoins ...
func (s *Transaction) MatchPlayerJoins(id, name string) (err error) {
	defer func() {
		s.errors = append(s.errors, err)
	}()

	m := s.inMemoryCopy.match(id)
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

// MatchPlayerLeaves ...
func (s *Transaction) MatchPlayerLeaves(id, name string) (err error) {
	defer func() {
		s.errors = append(s.errors, err)
	}()

	m := s.inMemoryCopy.match(id)
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
