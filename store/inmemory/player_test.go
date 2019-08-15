package inmemory

import (
	"plateau/protocol"
	"plateau/store"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlayerList(t *testing.T) {
	t.Parallel()

	s := &Store{}

	require.NoError(t, s.Open())
	defer func() {
		require.NoError(t, s.Close())
	}()

	trn, err := s.BeginTransaction()
	require.NoError(t, err)

	names, err := trn.PlayerList()
	require.NoError(t, err)
	require.Empty(t, names)

	require.NoError(t, trn.PlayerCreate(protocol.Player{Name: "foo"}))

	names, err = trn.PlayerList()
	require.NoError(t, err)
	require.Len(t, names, 1)
	require.Equal(t, "foo", names[0])
}

func TestPlayerCreate(t *testing.T) {
	t.Parallel()

	var (
		player = protocol.Player{Name: "foo"}

		s = &Store{}
	)

	require.NoError(t, s.Open())
	defer func() {
		require.NoError(t, s.Close())
	}()

	trn, err := s.BeginTransaction()
	require.NoError(t, err)

	require.NoError(t, trn.PlayerCreate(player))

	names, err := trn.PlayerList()
	require.NoError(t, err)
	require.Len(t, names, 1)
	require.Equal(t, "foo", names[0])

	require.IsType(t, store.DuplicateError(""), trn.PlayerCreate(player))
}

func TestPlayerRead(t *testing.T) {
	t.Parallel()

	var (
		player = protocol.Player{Name: "foo"}

		s = &Store{}
	)

	require.NoError(t, s.Open())
	defer func() {
		require.NoError(t, s.Close())
	}()

	trn, err := s.BeginTransaction()
	require.NoError(t, err)

	require.NoError(t, trn.PlayerCreate(player))

	p, err := trn.PlayerRead(player.Name)
	require.NoError(t, err)
	require.Equal(t, player.Name, p.Name)

	_, err = trn.PlayerRead("bar")
	require.IsType(t, store.DontExistError(""), err)
}

func TestPlayerIncreaseWins(t *testing.T) {
	t.Parallel()

	var (
		player = protocol.Player{Name: "foo"}

		s = &Store{}
	)

	require.NoError(t, s.Open())
	defer func() {
		require.NoError(t, s.Close())
	}()

	trn, err := s.BeginTransaction()
	require.NoError(t, err)

	require.NoError(t, trn.PlayerCreate(player))

	require.NoError(t, trn.PlayerIncreaseWins(player.Name, 2))

	p, err := trn.PlayerRead(player.Name)
	require.NoError(t, err)
	require.Equal(t, uint(2), p.Wins)

	require.IsType(t, store.DontExistError(""), trn.PlayerIncreaseWins("bar", 2))
}

func TestPlayerIncreaseLoses(t *testing.T) {
	t.Parallel()

	var (
		player = protocol.Player{Name: "foo"}

		s = &Store{}
	)

	require.NoError(t, s.Open())
	defer func() {
		require.NoError(t, s.Close())
	}()

	trn, err := s.BeginTransaction()
	require.NoError(t, err)

	require.NoError(t, trn.PlayerCreate(player))

	require.NoError(t, trn.PlayerIncreaseLoses(player.Name, 2))

	p, err := trn.PlayerRead(player.Name)
	require.NoError(t, err)
	require.Equal(t, uint(2), p.Loses)

	require.IsType(t, store.DontExistError(""), trn.PlayerIncreaseLoses("bar", 2))
}

func TestPlayerIncreaseTies(t *testing.T) {
	t.Parallel()

	var (
		player = protocol.Player{Name: "foo"}

		s = &Store{}
	)

	require.NoError(t, s.Open())
	defer func() {
		require.NoError(t, s.Close())
	}()

	trn, err := s.BeginTransaction()
	require.NoError(t, err)

	require.NoError(t, trn.PlayerCreate(player))

	require.NoError(t, trn.PlayerIncreaseTies(player.Name, 2))

	p, err := trn.PlayerRead(player.Name)
	require.NoError(t, err)
	require.Equal(t, uint(2), p.Ties)

	require.IsType(t, store.DontExistError(""), trn.PlayerIncreaseTies("bar", 2))
}
