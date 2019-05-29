package rethinkdb

import (
	"plateau/protocol"
	"testing"

	"github.com/stretchr/testify/require"
	rethinkdb "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

func TestPlayerConversion(t *testing.T) {
	t.Parallel()

	p := &player{}
	require.IsType(t, &protocol.Player{}, p.toProtocolStruct())
	require.IsType(t, p, playerFromProtocolStruct(p.toProtocolStruct()))
}

func TestPlayerList(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		str = &playerStore{mock}
	)

	mock.On(str.listTerm()).Return([]string{
		"foo", "bar",
	}, nil)

	players, err := str.List()
	require.NoError(t, err)
	require.Len(t, players, 2)

	mock.AssertExpectations(t)
}

func TestPlayerCreate(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		str = &playerStore{mock}

		givenName   = "foo"
		givenPlayer = protocol.Player{Name: givenName}
	)

	mock.On(str.createTerm(&givenPlayer)).Return(
		rethinkdb.WriteResponse{GeneratedKeys: []string{givenName}},
		nil)

	err := str.Create(givenPlayer)
	require.NoError(t, err)

	mock.AssertExpectations(t)
}

func TestPlayerRead(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		str = &playerStore{mock}

		givenPlayer    = player{Name: "foo"}
		expectedPlayer = *givenPlayer.toProtocolStruct()
	)

	mock.On(str.readTerm(givenPlayer.Name)).Return([]interface{}{
		givenPlayer,
	}, nil)

	player, err := str.Read(givenPlayer.Name)
	require.NoError(t, err)
	require.Equal(t, expectedPlayer.Name, player.Name)

	mock.AssertExpectations(t)
}

func TestPlayerIncreaseScore(t *testing.T) {
	t.Parallel()

	for _, field := range []string{"wins", "loses", "ties"} {
		var (
			mock = rethinkdb.NewMock()

			str = &playerStore{mock}

			givenPlayer = player{Name: "foo"}
		)

		mock.On(str.increaseScoreTerm(field, givenPlayer.Name, 1)).Return(rethinkdb.WriteResponse{}, nil)

		var err error

		switch field {
		case "wins":
			err = str.IncreaseWins(givenPlayer.Name, 1)
		case "loses":
			err = str.IncreaseLoses(givenPlayer.Name, 1)
		case "ties":
			err = str.IncreaseTies(givenPlayer.Name, 1)
		}

		require.NoError(t, err)

		mock.AssertExpectations(t)
	}
}
