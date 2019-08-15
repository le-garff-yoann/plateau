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

	NumberOfPlayersRequired uint
	Players                 map[string]protocol.Player

	Deals []deal
}

func matchFromProtocolStruct(m *protocol.Match) *match {
	var (
		players = make(map[string]protocol.Player)

		deals []deal
	)

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
		NumberOfPlayersRequired: m.NumberOfPlayersRequired,
		Players:                 players,
		Deals:                   deals,
	}
}

func (s *match) toProtocolStruct(pPlayers []*protocol.Player) *protocol.Match {
	var (
		players []protocol.Player

		pPlayersMap = make(map[string]protocol.Player)

		deals []protocol.Deal
	)

	for _, p := range pPlayers {
		pPlayersMap[p.Name] = *p
	}

	for pName, p := range s.Players {
		if pp, ok := pPlayersMap[pName]; ok {
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
		NumberOfPlayersRequired: s.NumberOfPlayersRequired,
		Players:                 players,
		Deals:                   deals,
	}
}

// MatchList implements the `store.Transaction` interface.
func (s *Transaction) MatchList() (IDs []string, err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}
	}()

	for _, m := range s.inMemoryCopy.Matchs {
		IDs = append(IDs, m.ID)
	}

	return IDs, nil
}

// MatchCreate implements the `store.Transaction` interface.
func (s *Transaction) MatchCreate(m protocol.Match) (_ string, err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}
	}()

	m.ID = uuid.NewV4().String()

	s.inMemoryCopy.Matchs = append(s.inMemoryCopy.Matchs, matchFromProtocolStruct(&m))

	return m.ID, nil
}

// MatchRead implements the `store.Transaction` interface.
func (s *Transaction) MatchRead(id string) (_ *protocol.Match, err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}
	}()

	m := s.inMemoryCopy.Match(id)
	if m == nil {
		return nil, store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	var matchCopy protocol.Match
	deepcopier.Copy(m.toProtocolStruct(s.inMemoryCopy.Players)).To(&matchCopy)

	return &matchCopy, nil
}

// MatchEndedAt implements the `store.Transaction` interface.
func (s *Transaction) MatchEndedAt(id string, val time.Time) (err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}

		s.matchNotifications = append(s.matchNotifications, store.MatchNotification{ID: id})
	}()

	m := s.inMemoryCopy.Match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	m.EndedAt = &val

	return nil
}

// MatchCreateDeal implements the `store.Transaction` interface.
func (s *Transaction) MatchCreateDeal(id string, deal protocol.Deal) (err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}

		s.matchNotifications = append(s.matchNotifications, store.MatchNotification{ID: id})
	}()
	m := s.inMemoryCopy.Match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	m.Deals = append(m.Deals, *dealFromProtocolStruct(&deal))

	return nil
}

// MatchUpdateCurrentDealHolder implements the `store.Transaction` interface.
func (s *Transaction) MatchUpdateCurrentDealHolder(id, newHolderName string) (err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}

		s.matchNotifications = append(s.matchNotifications, store.MatchNotification{ID: id})
	}()

	m := s.inMemoryCopy.Match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	m.Deals[len(m.Deals)-1].Holder.Name = newHolderName

	return nil
}

// MatchAddMessageToCurrentDeal implements the `store.Transaction` interface.
func (s *Transaction) MatchAddMessageToCurrentDeal(id string, message protocol.Message) (err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}

		s.matchNotifications = append(s.matchNotifications, store.MatchNotification{ID: id})
	}()

	m := s.inMemoryCopy.Match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	deal := &m.Deals[len(m.Deals)-1]
	deal.Messages = append(deal.Messages, message)

	return nil
}

// MatchPlayerJoins implements the `store.Transaction` interface.
func (s *Transaction) MatchPlayerJoins(id, name string) (err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}

		s.matchNotifications = append(s.matchNotifications, store.MatchNotification{ID: id})
	}()

	m := s.inMemoryCopy.Match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	if len(m.Players) >= int(m.NumberOfPlayersRequired) {
		return store.PlayerParticipationError(fmt.Sprintf(`There are no more spot in match %s`, id))
	}

	if _, ok := m.Players[name]; ok {
		return store.PlayerParticipationError(fmt.Sprintf(`The player "%s" is already in the match %s`, name, id))
	}

	m.Players[name] = protocol.Player{Name: name}

	return nil
}

// MatchPlayerLeaves implements the `store.Transaction` interface.
func (s *Transaction) MatchPlayerLeaves(id, name string) (err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}

		s.matchNotifications = append(s.matchNotifications, store.MatchNotification{ID: id})
	}()

	m := s.inMemoryCopy.Match(id)
	if m == nil {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	if _, ok := m.Players[name]; !ok {
		return store.PlayerParticipationError(fmt.Sprintf(`The player "%s" is already out of the match %s`, name, id))
	}

	delete(m.Players, name)

	return nil
}
