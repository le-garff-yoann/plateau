package rethinkdb

import (
	"plateau/store"
	"testing"

	"github.com/stretchr/testify/require"
	rethinkdb "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

func TestPlayerStoreConversion(t *testing.T) {
	t.Parallel()

	p := Player{}
	require.IsType(t, &store.Player{}, p.toStoreStruct())
	require.IsType(t, p, *playerFromStoreStruct(*p.toStoreStruct()))
}

func TestPlayerList(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		playerStore = &PlayerStore{mock}
	)

	mock.On(playerStore.listTerm()).Return([]string{
		"foo", "bar",
	}, nil)

	players, err := playerStore.List()
	require.NoError(t, err)
	require.Len(t, players, 2)

	mock.AssertExpectations(t)
}

func TestPlayerCreate(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		playerStore = &PlayerStore{mock}

		givenName   = "foo"
		givenPlayer = store.Player{Name: givenName}
	)

	mock.On(playerStore.createTerm(givenPlayer)).Return(
		rethinkdb.WriteResponse{GeneratedKeys: []string{givenName}},
		nil)

	err := playerStore.Create(givenPlayer)
	require.NoError(t, err)

	mock.AssertExpectations(t)
}

func TestPlayerRead(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		playerStore = &PlayerStore{mock}

		givenPlayer    = Player{Name: "foo"}
		expectedPlayer = *givenPlayer.toStoreStruct()
	)

	mock.On(playerStore.readTerm(givenPlayer.Name)).Return([]interface{}{
		givenPlayer,
	}, nil)

	player, err := playerStore.Read(givenPlayer.Name)
	require.NoError(t, err)
	require.IsType(t, &expectedPlayer, player)
	require.Equal(t, expectedPlayer.Name, player.Name)

	mock.AssertExpectations(t)
}

func TestPlayerIncreaseScore(t *testing.T) {
	t.Parallel()

	for _, field := range []string{"wins", "loses", "ties"} {
		var (
			mock = rethinkdb.NewMock()

			playerStore = &PlayerStore{mock}

			givenPlayer = Player{Name: "foo"}
		)

		mock.On(playerStore.increaseScoreTerm(field, givenPlayer.Name, 1)).Return(rethinkdb.WriteResponse{},
			nil)

		var err error

		switch field {
		case "wins":
			err = playerStore.IncreaseWins(givenPlayer.Name, 1)
		case "loses":
			err = playerStore.IncreaseLoses(givenPlayer.Name, 1)
		case "ties":
			err = playerStore.IncreaseTies(givenPlayer.Name, 1)
		}

		require.NoError(t, err)
	}
}
