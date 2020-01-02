package postgresql

import (
	"fmt"
	"plateau/protocol"
	"plateau/store"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	uuid "github.com/satori/go.uuid"
)

// Match ...
type Match struct {
	ID string

	CreatedAt time.Time
	EndedAt   *time.Time

	NumberOfPlayersRequired uint
	Players                 []Player `pg:"many2many:match_players"`

	Deals []Deal
}

// MatchPlayer ...
type MatchPlayer struct {
	MatchID  string `sql:",pk"`
	PlayerID string `sql:",pk"`
}

func matchFromProtocolStruct(m *protocol.Match) *Match {
	var (
		players []Player

		deals []Deal
	)

	for _, p := range m.Players {
		players = append(players, *playerFromProtocolStruct(&p))
	}

	for _, d := range m.Deals {
		deals = append(deals, *dealFromProtocolStruct(&d))
	}

	return &Match{
		ID:                      m.ID,
		CreatedAt:               m.CreatedAt,
		EndedAt:                 m.EndedAt,
		NumberOfPlayersRequired: m.NumberOfPlayersRequired,
		Players:                 players,
		Deals:                   deals,
	}
}

func (s *Match) toProtocolStruct() *protocol.Match {
	var (
		players []protocol.Player

		deals []protocol.Deal
	)

	for _, p := range s.Players {
		players = append(players, *p.toProtocolStruct())
	}

	for _, d := range s.Deals {
		deals = append(deals, *d.toProtocolStruct())
	}

	return &protocol.Match{
		ID:                      string(s.ID),
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

	var matchs []Match
	if err := s.tx.Model(&matchs).Select(); err != nil {
		return nil, err
	}

	for _, m := range matchs {
		IDs = append(IDs, m.toProtocolStruct().ID)
	}

	return IDs, err
}

// MatchCreate implements the `store.Transaction` interface.
func (s *Transaction) MatchCreate(m protocol.Match) (_ string, err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}
	}()

	m.ID = uuid.NewV4().String()

	match := matchFromProtocolStruct(&m)
	if err = s.tx.Insert(match); err != nil {
		return "", err
	}

	return match.toProtocolStruct().ID, nil
}

// MatchRead implements the `store.Transaction` interface.
func (s *Transaction) MatchRead(id string) (_ *protocol.Match, err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}
	}()

	match := Match{ID: matchFromProtocolStruct(&protocol.Match{ID: id}).ID}
	if err = s.tx.
		Model(&match).
		Relation("Players").
		WherePK().
		Select(); err != nil {
		if err == pg.ErrNoRows {
			return nil, store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
		}

		return nil, err
	}

	return match.toProtocolStruct(), nil
}

// MatchEndedAt implements the `store.Transaction` interface.
func (s *Transaction) MatchEndedAt(id string, val time.Time) (err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}

		s.matchNotifications = append(s.matchNotifications, store.MatchNotification{ID: id})
	}()

	var res orm.Result
	res, err = s.tx.
		Model(&Match{
			ID:      matchFromProtocolStruct(&protocol.Match{ID: id}).ID,
			EndedAt: &val,
		}).
		Column("ended_at").
		WherePK().
		Update()
	if err == nil && res.RowsAffected() == 0 {
		return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
	}

	return
}

// MatchCreateDeal implements the `store.Transaction` interface.
func (s *Transaction) MatchCreateDeal(id string, deal protocol.Deal) (err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}

		s.matchNotifications = append(s.matchNotifications, store.MatchNotification{ID: id})
	}()

	match := Match{ID: matchFromProtocolStruct(&protocol.Match{ID: id}).ID}
	if err = s.tx.
		Model(&match).
		Column("deals").
		WherePK().
		Select(); err != nil {
		if err == pg.ErrNoRows {
			return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
		}

		return err
	}

	match.Deals = append(match.Deals, *dealFromProtocolStruct(&deal))

	_, err = s.tx.
		Model(&match).
		Column("deals").
		WherePK().
		Update()

	return
}

// MatchUpdateCurrentDealHolder implements the `store.Transaction` interface.
func (s *Transaction) MatchUpdateCurrentDealHolder(id, newHolderName string) (err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}

		s.matchNotifications = append(s.matchNotifications, store.MatchNotification{ID: id})
	}()

	match := Match{ID: matchFromProtocolStruct(&protocol.Match{ID: id}).ID}
	if err = s.tx.
		Model(&match).
		Column("deals").
		WherePK().
		Select(); err != nil {
		if err == pg.ErrNoRows {
			return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
		}

		return err
	}

	match.Deals[len(match.Deals)-1].Holder = Player{ID: newHolderName}

	_, err = s.tx.
		Model(&match).
		Column("deals").
		WherePK().
		Update()

	return
}

// MatchAddMessageToCurrentDeal implements the `store.Transaction` interface.
func (s *Transaction) MatchAddMessageToCurrentDeal(id string, message protocol.Message) (err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}

		s.matchNotifications = append(s.matchNotifications, store.MatchNotification{ID: id})
	}()

	match := Match{ID: matchFromProtocolStruct(&protocol.Match{ID: id}).ID}
	if err = s.tx.
		Model(&match).
		Column("deals").
		WherePK().
		Select(); err != nil {
		if err == pg.ErrNoRows {
			return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
		}

		return err
	}

	match.Deals[len(match.Deals)-1].Messages = append(match.Deals[len(match.Deals)-1].Messages, *messageFromProtocolStruct(&message))

	_, err = s.tx.
		Model(&match).
		Column("deals").
		WherePK().
		Update()

	return
}

// MatchPlayerJoins implements the `store.Transaction` interface.
func (s *Transaction) MatchPlayerJoins(id, name string) (err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}

		s.matchNotifications = append(s.matchNotifications, store.MatchNotification{ID: id})
	}()

	match := matchFromProtocolStruct(&protocol.Match{ID: id})

	err = s.tx.
		Model(match).
		WherePK().
		Column("number_of_players_required").
		Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return store.DontExistError(fmt.Sprintf(`The match %s doesn't exist`, id))
		}

		return
	}

	playerCount, err := s.tx.
		Model(&MatchPlayer{}).
		Where("match_id = ?", match.ID).
		Count()
	if err != nil {
		return err
	}

	if uint(playerCount) >= match.NumberOfPlayersRequired {
		return store.PlayerParticipationError(fmt.Sprintf(`There are no more spot in match %s`, id))
	}

	if err = s.tx.Insert(&MatchPlayer{
		MatchID:  match.ID,
		PlayerID: name,
	}); err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() {
			return store.PlayerParticipationError(fmt.Sprintf(`The player "%s" is already in the match %s`, name, id))
		}
	}

	return
}

// MatchPlayerLeaves implements the `store.Transaction` interface.
func (s *Transaction) MatchPlayerLeaves(id, name string) (err error) {
	defer func() {
		if err != nil {
			s.errors = append(s.errors, err)
		}

		s.matchNotifications = append(s.matchNotifications, store.MatchNotification{ID: id})
	}()

	if err = s.tx.Delete(&MatchPlayer{
		MatchID:  matchFromProtocolStruct(&protocol.Match{ID: id}).ID,
		PlayerID: name,
	}); err == pg.ErrNoRows {
		return store.PlayerParticipationError(fmt.Sprintf(`The player "%s" is already out of the match %s`, name, id))
	}

	return
}
