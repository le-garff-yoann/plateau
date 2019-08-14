package inmemory

import (
	"fmt"
	"plateau/protocol"
	"plateau/store"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestMatchCreateAndList(t *testing.T) {
	t.Parallel()

	s := &Store{}

	require.NoError(t, s.Open())
	defer func() {
		require.NoError(t, s.Close())
	}()

	trn := s.BeginTransaction()

	IDs, err := trn.MatchList()
	require.NoError(t, err)
	require.Empty(t, IDs)

	_, err = trn.MatchCreate(protocol.Match{})
	require.NoError(t, err)

	IDs, err = trn.MatchList()
	require.NoError(t, err)
	require.Len(t, IDs, 1)
	require.NotEmpty(t, IDs[0])
}

func TestMatchRead(t *testing.T) {
	t.Parallel()

	s := &Store{}

	require.NoError(t, s.Open())
	defer func() {
		require.NoError(t, s.Close())
	}()

	trn := s.BeginTransaction()

	id, err := trn.MatchCreate(protocol.Match{})
	require.NoError(t, err)

	_, err = uuid.FromString(id)
	require.NoError(t, err)

	m, err := trn.MatchRead(id)
	require.NoError(t, err)
	require.Equal(t, id, m.ID)

	_, err = trn.PlayerRead(fmt.Sprintf("%s0", id))
	require.IsType(t, store.DontExistError(""), err)
}

func TestMatchEndedAt(t *testing.T) {
	t.Parallel()

	s := &Store{}

	require.NoError(t, s.Open())
	defer func() {
		require.NoError(t, s.Close())
	}()

	trn := s.BeginTransaction()

	id, err := trn.MatchCreate(protocol.Match{})
	require.NoError(t, err)

	err = trn.MatchEndedAt(id, time.Now())
	require.NoError(t, err)

	m, err := trn.MatchRead(id)
	require.NoError(t, err)
	require.NotNil(t, id, m.EndedAt)

	require.IsType(t, store.DontExistError(""),
		trn.MatchEndedAt(fmt.Sprintf("%s0", id), time.Now()))
}

func TestMatchCreateDeal(t *testing.T) {
	t.Parallel()

	var (
		s = &Store{}

		deal = protocol.Deal{Holder: protocol.Player{Name: "foo"}}
	)

	require.NoError(t, s.Open())
	defer func() {
		require.NoError(t, s.Close())
	}()

	trn := s.BeginTransaction()

	id, err := trn.MatchCreate(protocol.Match{})
	require.NoError(t, err)

	err = trn.MatchCreateDeal(id, deal)
	require.NoError(t, err)

	m, err := trn.MatchRead(id)
	require.NoError(t, err)
	require.Len(t, m.Deals, 1)
	require.Equal(t, deal.Holder.Name, m.Deals[0].Holder.Name)

	require.IsType(t, store.DontExistError(""),
		trn.MatchEndedAt(fmt.Sprintf("%s0", id), time.Now()))
}

func TestMatchUpdateCurrentDealHolder(t *testing.T) {
	t.Parallel()

	var (
		s = &Store{}

		deal = protocol.Deal{Holder: protocol.Player{Name: "foo", Wins: 1}}
	)

	require.NoError(t, s.Open())
	defer func() {
		require.NoError(t, s.Close())
	}()

	trn := s.BeginTransaction()

	id, err := trn.MatchCreate(protocol.Match{})
	require.NoError(t, err)

	require.NoError(t, trn.MatchCreateDeal(id, deal))

	err = trn.MatchUpdateCurrentDealHolder(id, "bar")
	require.NoError(t, err)

	m, err := trn.MatchRead(id)
	require.NoError(t, err)
	require.Len(t, m.Deals, 1)
	require.Equal(t, "bar", m.Deals[0].Holder.Name)

	err = trn.MatchUpdateCurrentDealHolder(id, "foo")
	require.NoError(t, err)

	m, err = trn.MatchRead(id)
	require.NoError(t, err)
	require.Len(t, m.Deals, 1)
	require.Equal(t, "foo", m.Deals[0].Holder.Name)
	require.Equal(t, uint(1), m.Deals[0].Holder.Wins)

	require.IsType(t, store.DontExistError(""),
		trn.MatchEndedAt(fmt.Sprintf("%s0", id), time.Now()))
}

func TestMatchPlayerJoins(t *testing.T) {
	t.Parallel()

	s := &Store{}

	require.NoError(t, s.Open())
	defer func() {
		require.NoError(t, s.Close())
	}()

	trn := s.BeginTransaction()

	id, err := trn.MatchCreate(protocol.Match{NumberOfPlayersRequired: 2})
	require.NoError(t, err)

	require.NoError(t, trn.MatchPlayerJoins(id, "foo"))

	m, err := trn.MatchRead(id)
	require.NoError(t, err)
	require.Equal(t, "foo", m.Players[0].Name)

	require.IsType(t, store.PlayerParticipationError(""),
		trn.MatchPlayerJoins(id, "foo"))

	require.NoError(t, trn.MatchPlayerJoins(id, "bar"))
	require.IsType(t, store.PlayerParticipationError(""),
		trn.MatchPlayerJoins(id, "baz"))
}

func TestMatchPlayerLeaves(t *testing.T) {
	t.Parallel()

	s := &Store{}

	require.NoError(t, s.Open())
	defer func() {
		require.NoError(t, s.Close())
	}()

	trn := s.BeginTransaction()

	id, err := trn.MatchCreate(protocol.Match{NumberOfPlayersRequired: 1})
	require.NoError(t, err)

	require.NoError(t, trn.MatchPlayerJoins(id, "foo"))

	require.NoError(t, trn.MatchPlayerLeaves(id, "foo"))

	m, err := trn.MatchRead(id)
	require.NoError(t, err)
	require.Empty(t, m.Players)

	require.IsType(t, store.PlayerParticipationError(""),
		trn.MatchPlayerLeaves(id, "foo"))
}
