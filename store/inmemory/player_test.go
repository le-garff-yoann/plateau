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
	s.Open()

	names, err := s.Players().List()
	require.NoError(t, err)
	require.Empty(t, names)

	s.Players().Create(protocol.Player{Name: "foo"})

	names, err = s.Players().List()
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

	s.Open()

	require.NoError(t, s.Players().Create(player))

	names, err := s.Players().List()
	require.NoError(t, err)
	require.Len(t, names, 1)
	require.Equal(t, "foo", names[0])

	require.IsType(t, store.DuplicateError(""), s.Players().Create(player))
}

func TestPlayerRead(t *testing.T) {
	t.Parallel()

	var (
		player = protocol.Player{Name: "foo"}

		s = &Store{}
	)

	s.Open()

	s.Players().Create(player)

	p, err := s.Players().Read(player.Name)
	require.NoError(t, err)
	require.Equal(t, player.Name, p.Name)

	_, err = s.Players().Read("bar")
	require.IsType(t, store.DontExistError(""), err)
}

func TestPlayerIncreaseWins(t *testing.T) {
	t.Parallel()

	var (
		player = protocol.Player{Name: "foo"}

		s = &Store{}
	)

	s.Open()

	s.Players().Create(player)

	require.NoError(t, s.Players().IncreaseWins(player.Name, 2))

	p, _ := s.Players().Read(player.Name)
	require.Equal(t, uint(2), p.Wins)

	require.IsType(t, store.DontExistError(""), s.Players().IncreaseWins("bar", 2))
}

func TestPlayerIncreaseLoses(t *testing.T) {
	t.Parallel()

	var (
		player = protocol.Player{Name: "foo"}

		s = &Store{}
	)

	s.Open()

	s.Players().Create(player)

	require.NoError(t, s.Players().IncreaseLoses(player.Name, 2))

	p, _ := s.Players().Read(player.Name)
	require.Equal(t, uint(2), p.Loses)

	require.IsType(t, store.DontExistError(""), s.Players().IncreaseLoses("bar", 2))
}

func TestPlayerIncreaseTies(t *testing.T) {
	t.Parallel()

	var (
		player = protocol.Player{Name: "foo"}

		s = &Store{}
	)

	s.Open()

	s.Players().Create(player)

	require.NoError(t, s.Players().IncreaseTies(player.Name, 2))

	p, _ := s.Players().Read(player.Name)
	require.Equal(t, uint(2), p.Ties)

	require.IsType(t, store.DontExistError(""), s.Players().IncreaseTies("bar", 2))
}
