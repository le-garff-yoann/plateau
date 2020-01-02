package store

import (
	"fmt"
	"plateau/protocol"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

// TestStore is accompanied by a list of methods that act
// as an agnostic test suite for the implementation of `Store`.
//	- *SetupCb* is optionnal and is called after every `WrapTest()` call.
//	- *TeardownCb* is optionnal and is called before every `WrapTest()` call.
type TestStore struct {
	T *testing.T

	Str                 Store
	SetupCb, TeardownCb func(testStr *TestStore)
}

// WrapTest calls `Str.Open()`, *cb* and `Str.Close()`.
func (s *TestStore) WrapTest(cb func(*TestStore)) {
	require.NoError(s.T, s.Str.Open())
	if s.SetupCb != nil {
		s.SetupCb(s)
	}

	defer func() {
		if s.TeardownCb != nil {
			s.TeardownCb(s)
		}

		require.NoError(s.T, s.Str.Close())
	}()

	cb(s)
}

// WrapTransaction starts a `Transaction`, calls *cb* and
// `Abort()` the transaction within a `WrapTest()` call.
func (s *TestStore) WrapTransaction(cb func(*testing.T, Transaction)) {
	s.WrapTest(func(s *TestStore) {
		trn, err := s.Str.BeginTransaction()
		require.NoError(s.T, err)

		defer func() {
			require.NoError(s.T, trn.Abort())
		}()

		cb(s.T, trn)
	})
}

// TestNotificationsChannel runs the test suite
// related to the subject in the method name.
func (s *TestStore) TestNotificationsChannel() {
	s.WrapTest(func(s *TestStore) {
		ch := make(chan interface{})

		require.NoError(s.T, s.Str.RegisterNotificationsChannel(ch))
		require.NoError(s.T, s.Str.UnregisterNotificationsChannel(ch))
	})
}

// TestBeginTransactionCommit runs the test suite
// related to the subject in the method name.
func (s *TestStore) TestBeginTransactionCommit() {
	s.WrapTest(func(s *TestStore) {
		trn, err := s.Str.BeginTransaction()
		require.NoError(s.T, err)

		id, err := trn.MatchCreate(protocol.Match{})
		require.NoError(s.T, err)

		_, err = trn.MatchRead(id)
		require.NoError(s.T, err)

		require.False(s.T, trn.Closed())
		require.NotPanics(s.T, func() {
			require.NoError(s.T, trn.Commit())
		})
		require.Panics(s.T, func() {
			require.NoError(s.T, trn.Commit())
		})
		require.Panics(s.T, func() {
			require.NoError(s.T, trn.Abort())
		})

		require.True(s.T, trn.Closed())

		trn, err = s.Str.BeginTransaction()
		require.NoError(s.T, err)

		_, err = trn.MatchRead(id)
		require.NoError(s.T, err)
		require.NotPanics(s.T, func() {
			require.NoError(s.T, trn.Commit())
		})
	})
}

// TestBeginTransactionAbort runs the test suite
// related to the subject in the method name.
func (s *TestStore) TestBeginTransactionAbort() {
	s.WrapTest(func(s *TestStore) {
		trn, err := s.Str.BeginTransaction()
		require.NoError(s.T, err)

		id, err := trn.MatchCreate(protocol.Match{})
		require.NoError(s.T, err)

		_, err = trn.MatchRead(id)
		require.NoError(s.T, err)

		require.False(s.T, trn.Closed())
		require.NoError(s.T, trn.Abort())
		require.Panics(s.T, func() {
			require.NoError(s.T, trn.Commit())
		})
		require.Panics(s.T, func() {
			require.NoError(s.T, trn.Abort())
		})

		require.True(s.T, trn.Closed())

		require.Empty(s.T, trn.Errors())

		trn, err = s.Str.BeginTransaction()
		require.NoError(s.T, err)

		_, err = trn.MatchRead(id)
		require.IsType(s.T, DontExistError(""), err)
	})
}

// TestPlayerList runs the test suite
// related to the subject in the method name.
func (s *TestStore) TestPlayerList() {
	s.WrapTransaction(func(t *testing.T, trn Transaction) {
		names, err := trn.PlayerList()
		require.NoError(t, err)
		require.Empty(t, names)

		require.NoError(t, trn.PlayerCreate(protocol.Player{Name: "foo"}))

		names, err = trn.PlayerList()
		require.NoError(t, err)
		require.Len(t, names, 1)
		require.Equal(t, "foo", names[0])
	})
}

// TestPlayerCreate runs the test suite
// related to the subject in the method name.
func (s *TestStore) TestPlayerCreate() {
	s.WrapTransaction(func(t *testing.T, trn Transaction) {
		player := protocol.Player{Name: "foo"}

		require.NoError(t, trn.PlayerCreate(player))

		names, err := trn.PlayerList()
		require.NoError(t, err)
		require.Len(t, names, 1)
		require.Equal(t, "foo", names[0])

		require.IsType(t, DuplicateError(""), trn.PlayerCreate(player))
	})
}

// TestPlayerRead runs the test suite
// related to the subject in the method name.
func (s *TestStore) TestPlayerRead() {
	s.WrapTransaction(func(t *testing.T, trn Transaction) {
		player := protocol.Player{Name: "foo"}

		require.NoError(t, trn.PlayerCreate(player))

		p, err := trn.PlayerRead(player.Name)
		require.NoError(t, err)
		require.Equal(t, player.Name, p.Name)

		_, err = trn.PlayerRead("bar")
		require.IsType(t, DontExistError(""), err)
	})
}

// TestPlayerIncreaseWins runs the test suite
// related to the subject in the method name.
func (s *TestStore) TestPlayerIncreaseWins() {
	s.WrapTransaction(func(t *testing.T, trn Transaction) {
		player := protocol.Player{Name: "foo"}

		require.NoError(t, trn.PlayerCreate(player))

		require.NoError(t, trn.PlayerIncreaseWins(player.Name, 2))

		p, err := trn.PlayerRead(player.Name)
		require.NoError(t, err)
		require.Equal(t, uint(2), p.Wins)

		require.IsType(t, DontExistError(""), trn.PlayerIncreaseWins("bar", 2))
	})
}

// TestPlayerIncreaseLoses runs the test suite
// related to the subject in the method name.
func (s *TestStore) TestPlayerIncreaseLoses() {
	s.WrapTransaction(func(t *testing.T, trn Transaction) {
		player := protocol.Player{Name: "foo"}

		require.NoError(t, trn.PlayerCreate(player))

		require.NoError(t, trn.PlayerIncreaseLoses(player.Name, 2))

		p, err := trn.PlayerRead(player.Name)
		require.NoError(t, err)
		require.Equal(t, uint(2), p.Loses)

		require.IsType(t, DontExistError(""), trn.PlayerIncreaseLoses("bar", 2))
	})
}

// TestPlayerIncreaseTies runs the test suite
// related to the subject in the method name.
func (s *TestStore) TestPlayerIncreaseTies() {
	s.WrapTransaction(func(t *testing.T, trn Transaction) {
		player := protocol.Player{Name: "foo"}

		require.NoError(t, trn.PlayerCreate(player))

		require.NoError(t, trn.PlayerIncreaseTies(player.Name, 2))

		p, err := trn.PlayerRead(player.Name)
		require.NoError(t, err)
		require.Equal(t, uint(2), p.Ties)

		require.IsType(t, DontExistError(""), trn.PlayerIncreaseTies("bar", 2))
	})
}

// TestMatchCreateAndList runs the test suite
// related to the subject in the method name.
func (s *TestStore) TestMatchCreateAndList() {
	s.WrapTransaction(func(t *testing.T, trn Transaction) {
		IDs, err := trn.MatchList()
		require.NoError(t, err)
		require.Empty(t, IDs)

		_, err = trn.MatchCreate(protocol.Match{})
		require.NoError(t, err)

		IDs, err = trn.MatchList()
		require.NoError(t, err)
		require.Len(t, IDs, 1)
		require.NotEmpty(t, IDs[0])
	})
}

// TestMatchRead runs the test suite
// related to the subject in the method name.
func (s *TestStore) TestMatchRead() {
	s.WrapTransaction(func(t *testing.T, trn Transaction) {
		id, err := trn.MatchCreate(protocol.Match{})
		require.NoError(t, err)

		_, err = uuid.FromString(id)
		require.NoError(t, err)

		m, err := trn.MatchRead(id)
		require.NoError(t, err)
		require.Equal(t, id, m.ID)

		_, err = trn.PlayerRead(fmt.Sprintf("%s0", id))
		require.IsType(t, DontExistError(""), err)
	})
}

// TestMatchEndedAt runs the test suite
// related to the subject in the method name.
func (s *TestStore) TestMatchEndedAt() {
	s.WrapTransaction(func(t *testing.T, trn Transaction) {
		id, err := trn.MatchCreate(protocol.Match{})
		require.NoError(t, err)

		err = trn.MatchEndedAt(id, time.Now())
		require.NoError(t, err)

		m, err := trn.MatchRead(id)
		require.NoError(t, err)
		require.NotNil(t, id, m.EndedAt)

		require.IsType(t, DontExistError(""),
			trn.MatchEndedAt(fmt.Sprintf("%s0", id), time.Now()))
	})
}

// TestMatchCreateDeal runs the test suite
// related to the subject in the method name.
func (s *TestStore) TestMatchCreateDeal() {
	s.WrapTransaction(func(t *testing.T, trn Transaction) {
		deal := protocol.Deal{Holder: protocol.Player{Name: "foo"}}

		id, err := trn.MatchCreate(protocol.Match{})
		require.NoError(t, err)

		err = trn.MatchCreateDeal(id, deal)
		require.NoError(t, err)

		m, err := trn.MatchRead(id)
		require.NoError(t, err)
		require.Len(t, m.Deals, 1)
		require.Equal(t, deal.Holder.Name, m.Deals[0].Holder.Name)
	})
}

// TestMatchUpdateCurrentDealHolder runs the test suite
// related to the subject in the method name.
func (s *TestStore) TestMatchUpdateCurrentDealHolder() {
	s.WrapTransaction(func(t *testing.T, trn Transaction) {
		deal := protocol.Deal{Holder: protocol.Player{Name: "foo", Wins: 1}}

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
	})
}

// TestMatchPlayerJoins runs the test suite
// related to the subject in the method name.
func (s *TestStore) TestMatchPlayerJoins() {
	var (
		playerA = protocol.Player{Name: "foo"}
		playerB = protocol.Player{Name: "bar"}
	)

	baseTests := func(t *testing.T, trn Transaction) string {
		id, err := trn.MatchCreate(protocol.Match{NumberOfPlayersRequired: 2})
		require.NoError(t, err)

		require.NoError(t, trn.PlayerCreate(playerA))
		require.NoError(t, trn.PlayerCreate(playerB))

		require.NoError(t, trn.MatchPlayerJoins(id, playerA.Name))

		m, err := trn.MatchRead(id)
		require.NoError(t, err)
		require.Equal(t, playerA.Name, m.Players[0].Name)

		return id
	}

	s.WrapTransaction(func(t *testing.T, trn Transaction) {
		id := baseTests(t, trn)

		require.IsType(t, PlayerParticipationError(""),
			trn.MatchPlayerJoins(id, playerA.Name))
	})

	s.WrapTransaction(func(t *testing.T, trn Transaction) {
		id := baseTests(t, trn)

		require.NoError(t, trn.MatchPlayerJoins(id, playerB.Name))
		require.IsType(t, PlayerParticipationError(""),
			trn.MatchPlayerJoins(id, "baz"))
	})
}

// TestMatchPlayerLeaves runs the test suite
// related to the subject in the method name.
func (s *TestStore) TestMatchPlayerLeaves() {
	s.WrapTransaction(func(t *testing.T, trn Transaction) {
		id, err := trn.MatchCreate(protocol.Match{NumberOfPlayersRequired: 1})
		require.NoError(t, err)

		player := protocol.Player{Name: "foo"}
		require.NoError(t, trn.PlayerCreate(player))

		require.NoError(t, trn.MatchPlayerJoins(id, player.Name))

		require.NoError(t, trn.MatchPlayerLeaves(id, player.Name))

		m, err := trn.MatchRead(id)
		require.NoError(t, err)
		require.Empty(t, m.Players)

		require.IsType(t, PlayerParticipationError(""),
			trn.MatchPlayerLeaves(id, player.Name))
	})
}
