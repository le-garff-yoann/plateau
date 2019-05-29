package rethinkdb

import (
	"plateau/protocol"
	"plateau/store"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	rethinkdb "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

func TestMatchConversion(t *testing.T) {
	t.Parallel()

	m := &match{}
	require.IsType(t, &protocol.Match{}, m.toProtocolStruct())
	require.IsType(t, m, matchFromProtocolStruct(m.toProtocolStruct()))
}

func TestMatchList(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		str = &matchStore{mock}
	)

	mock.On(str.listTerm()).Return([]string{
		"foo", "bar",
	}, nil)

	matchs, err := str.List()
	require.NoError(t, err)
	require.Len(t, matchs, 2)

	mock.AssertExpectations(t)
}

func TestMatchCreate(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		str = &matchStore{mock}

		givenMatch = protocol.Match{ID: "foo"}
	)

	mock.On(str.createTerm(&givenMatch)).Return(
		rethinkdb.WriteResponse{GeneratedKeys: []string{givenMatch.ID}},
		nil)

	id, err := str.Create(givenMatch)
	require.NoError(t, err)
	require.Equal(t, givenMatch.ID, id)

	mock.AssertExpectations(t)
}

func TestMatchRead(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		str = &matchStore{mock}

		givenID       = "foo"
		givenMatch    = match{ID: givenID}
		expectedMatch = *givenMatch.toProtocolStruct()
	)

	mock.On(str.readTerm(givenID)).Return([]interface{}{
		givenMatch,
	}, nil)

	match, err := str.Read(givenID)
	require.NoError(t, err)
	require.Equal(t, expectedMatch.ID, match.ID)

	mock.AssertExpectations(t)
}

func TestMatchEndedAt(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		str = &matchStore{mock}

		givenMatch  = match{ID: "foo"}
		endedAtTime = time.Now()
	)

	mock.On(str.endedAtTerm(givenMatch.ID, &endedAtTime)).Return(rethinkdb.WriteResponse{}, nil)

	err := str.EndedAt(givenMatch.ID, endedAtTime)
	require.NoError(t, err)

	mock.AssertExpectations(t)
}

func TestMatchConnectPlayer(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		str = &matchStore{mock}

		givenID         = "foo"
		givenPlayerName = "bar"
	)

	mock.
		On(str.connectPlayerTerm(givenID, givenPlayerName)).
		Return(rethinkdb.WriteResponse{Replaced: 1}, nil).
		Once()

	err := str.ConnectPlayer(givenID, givenPlayerName)
	require.NoError(t, err)

	mock.
		On(str.connectPlayerTerm(givenID, givenPlayerName)).
		Return(rethinkdb.WriteResponse{Replaced: 0}, nil).
		Once()

	err = str.ConnectPlayer(givenID, givenPlayerName)
	require.Error(t, err)

	mock.AssertExpectations(t)
}

func TestMatchDisconnectPlayer(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		str = &matchStore{mock}

		givenID         = "foo"
		givenPlayerName = "bar"
	)

	mock.On(str.disconnectPlayerTerm(givenID, givenPlayerName)).Return(rethinkdb.WriteResponse{}, nil)

	err := str.DisconnectPlayer(givenID, givenPlayerName)
	require.NoError(t, err)

	mock.AssertExpectations(t)
}

func TestMatchPlayerJoins(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		str = &matchStore{mock}

		givenID         = "foo"
		givenPlayerName = "bar"
	)

	mock.
		On(str.playerJoinsTerm(givenID, givenPlayerName)).
		Return(rethinkdb.WriteResponse{Replaced: 1}, nil).
		Once()

	err := str.PlayerJoins(givenID, givenPlayerName)
	require.NoError(t, err)

	mock.
		On(str.playerJoinsTerm(givenID, givenPlayerName)).
		Return(rethinkdb.WriteResponse{Replaced: 0}, nil).
		Once()

	err = str.PlayerJoins(givenID, givenPlayerName)
	require.Error(t, err)

	mock.AssertExpectations(t)
}

func TestMatchPlayerLeaves(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		str = &matchStore{mock}

		givenID         = "foo"
		givenPlayerName = "bar"
	)

	mock.
		On(str.playerLeavesTerm(givenID, givenPlayerName)).
		Return(rethinkdb.WriteResponse{Replaced: 1}, nil).
		Once()

	err := str.PlayerLeaves(givenID, givenPlayerName)
	require.NoError(t, err)

	mock.
		On(str.playerLeavesTerm(givenID, givenPlayerName)).
		Return(rethinkdb.WriteResponse{Replaced: 0}, nil).
		Once()

	err = str.PlayerLeaves(givenID, givenPlayerName)
	require.Error(t, err)

	mock.AssertExpectations(t)
}

func TestMatchCreateTransaction(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		str = &matchStore{mock}

		givenID  = "foo"
		givenTrx = protocol.Transaction{}
	)

	mock.
		On(str.createTransactionTerm(givenID, &givenTrx)).
		Return(rethinkdb.WriteResponse{}, nil)

	err := str.CreateTransaction(givenID, givenTrx)
	require.NoError(t, err)

	mock.AssertExpectations(t)
}

func TestMatchUpdateCurrentTransactionHolder(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		str        = &matchStore{mock}
		givenID    = "foo"
		holderName = "bar"
	)

	mock.
		On(str.updateCurrentTransactionHolderTerm(givenID, holderName)).
		Return(rethinkdb.WriteResponse{Replaced: 1}, nil)

	err := str.UpdateCurrentTransactionHolder(givenID, holderName)
	require.NoError(t, err)

	mock.AssertExpectations(t)
}

func TestMatchAddMessageToCurrentTransaction(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		str     = &matchStore{mock}
		givenID = "foo"
		msg     = protocol.Message{}
	)

	mock.
		On(str.addMessageToCurrentTransaction(givenID, &msg)).
		Return(rethinkdb.WriteResponse{Replaced: 1}, nil)

	err := str.AddMessageToCurrentTransaction(givenID, msg)
	require.NoError(t, err)

	mock.AssertExpectations(t)
}

func TestMatchCreateTransactionsChangeIterator(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		str = &matchStore{mock}

		givenID                  = "foo"
		givenMatchChangeResponse = matchChangeResponse{
			OldValue: match{Transactions: []transaction{}},
			NewValue: match{Transactions: []transaction{transaction{}, transaction{}}},
		}
	)

	mock.On(str.matchChangesTerm(givenID)).Return(givenMatchChangeResponse, nil)

	iterator, err := str.CreateTransactionsChangeIterator(givenID)
	require.NoError(t, err)

	var trxChange store.TransactionChange

	for range givenMatchChangeResponse.NewValue.Transactions {
		iterator.Next(&trxChange)

		require.Nil(t, trxChange.Old)
		require.NotNil(t, trxChange.New)
	}

	mock.AssertExpectations(t)
}
