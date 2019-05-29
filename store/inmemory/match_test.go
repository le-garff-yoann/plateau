package inmemory

import (
	"fmt"
	"plateau/protocol"
	"plateau/store"
	"sync"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestMatchStoreConversion(t *testing.T) {
	t.Parallel()

	var (
		m = &match{}

		players = []*protocol.Player{}
	)

	require.IsType(t, &protocol.Match{}, m.toProtocolStruct(players))
	require.IsType(t, m, matchFromProtocolStruct(m.toProtocolStruct(players)))
}

func TestMatchCreateAndList(t *testing.T) {
	t.Parallel()

	s := &Store{}
	s.Open()

	IDs, err := s.Matchs().List()
	require.NoError(t, err)
	require.Empty(t, IDs)

	_, err = s.Matchs().Create(protocol.Match{})
	require.NoError(t, err)

	IDs, err = s.Matchs().List()
	require.NoError(t, err)
	require.Len(t, IDs, 1)
	require.NotEmpty(t, IDs[0])
}

func TestMatchRead(t *testing.T) {
	t.Parallel()

	s := &Store{}
	s.Open()

	id, _ := s.Matchs().Create(protocol.Match{})

	_, err := uuid.FromString(id)
	require.NoError(t, err)

	m, err := s.Matchs().Read(id)
	require.NoError(t, err)
	require.Equal(t, id, m.ID)

	_, err = s.Players().Read(fmt.Sprintf("%s0", id))
	require.IsType(t, store.DontExistError(""), err)
}

func TestMatchEndedAt(t *testing.T) {
	t.Parallel()

	s := &Store{}
	s.Open()

	id, _ := s.Matchs().Create(protocol.Match{})

	err := s.Matchs().EndedAt(id, time.Now())
	require.NoError(t, err)

	m, _ := s.Matchs().Read(id)
	require.NotNil(t, id, m.EndedAt)

	require.IsType(t, store.DontExistError(""),
		s.Matchs().EndedAt(fmt.Sprintf("%s0", id), time.Now()))
}

func TestMatchCreateTransaction(t *testing.T) {
	t.Parallel()

	var (
		s = &Store{}

		trx = protocol.Transaction{Holder: protocol.Player{Name: "foo"}}
	)

	s.Open()

	id, _ := s.Matchs().Create(protocol.Match{})

	err := s.Matchs().CreateTransaction(id, trx)
	require.NoError(t, err)

	m, _ := s.Matchs().Read(id)
	require.Len(t, m.Transactions, 1)
	require.Equal(t, trx.Holder.Name, m.Transactions[0].Holder.Name)

	require.IsType(t, store.DontExistError(""),
		s.Matchs().EndedAt(fmt.Sprintf("%s0", id), time.Now()))
}

func TestMatchUpdateCurrentTransactionHolder(t *testing.T) {
	t.Parallel()

	var (
		s = &Store{}

		trx = protocol.Transaction{Holder: protocol.Player{Name: "foo", Wins: 1}}
	)

	s.Open()

	id, _ := s.Matchs().Create(protocol.Match{})

	s.Matchs().CreateTransaction(id, trx)

	err := s.Matchs().UpdateCurrentTransactionHolder(id, "bar")
	require.NoError(t, err)

	m, _ := s.Matchs().Read(id)
	require.Len(t, m.Transactions, 1)
	require.Equal(t, "bar", m.Transactions[0].Holder.Name)

	err = s.Matchs().UpdateCurrentTransactionHolder(id, "foo")
	require.NoError(t, err)

	m, _ = s.Matchs().Read(id)
	require.Len(t, m.Transactions, 1)
	require.Equal(t, "foo", m.Transactions[0].Holder.Name)
	require.Equal(t, uint(1), m.Transactions[0].Holder.Wins)

	require.IsType(t, store.DontExistError(""),
		s.Matchs().EndedAt(fmt.Sprintf("%s0", id), time.Now()))
}

func TestMatchConnectPlayer(t *testing.T) {
	t.Parallel()

	s := &Store{}
	s.Open()

	id, _ := s.Matchs().Create(protocol.Match{})

	require.NoError(t, s.Matchs().ConnectPlayer(id, "foo"))

	m, _ := s.Matchs().Read(id)
	require.Equal(t, "foo", m.ConnectedPlayers[0].Name)

	require.IsType(t, store.PlayerConnectionError(""),
		s.Matchs().ConnectPlayer(id, "foo"))
}

func TestMatchDisconnectPlayer(t *testing.T) {
	t.Parallel()

	s := &Store{}
	s.Open()

	id, _ := s.Matchs().Create(protocol.Match{})

	s.Matchs().ConnectPlayer(id, "foo")

	require.NoError(t, s.Matchs().DisconnectPlayer(id, "foo"))

	m, _ := s.Matchs().Read(id)
	require.Empty(t, m.ConnectedPlayers)
}

func TestMatchPlayerJoins(t *testing.T) {
	t.Parallel()

	s := &Store{}
	s.Open()

	id, _ := s.Matchs().Create(protocol.Match{NumberOfPlayersRequired: 2})

	require.NoError(t, s.Matchs().PlayerJoins(id, "foo"))

	m, _ := s.Matchs().Read(id)
	require.Equal(t, "foo", m.Players[0].Name)

	require.IsType(t, store.PlayerParticipationError(""),
		s.Matchs().PlayerJoins(id, "foo"))

	require.NoError(t, s.Matchs().PlayerJoins(id, "bar"))
	require.IsType(t, store.PlayerParticipationError(""),
		s.Matchs().PlayerJoins(id, "baz"))
}

func TestMatchPlayerLeaves(t *testing.T) {
	t.Parallel()

	s := &Store{}
	s.Open()

	id, _ := s.Matchs().Create(protocol.Match{NumberOfPlayersRequired: 1})

	require.NoError(t, s.Matchs().PlayerJoins(id, "foo"))

	require.NoError(t, s.Matchs().PlayerLeaves(id, "foo"))

	m, _ := s.Matchs().Read(id)
	require.Empty(t, m.Players)

	require.IsType(t, store.PlayerParticipationError(""),
		s.Matchs().PlayerLeaves(id, "foo"))
}

func TestMatchCreateTransactionsChangeIterator(t *testing.T) {
	t.Parallel()

	s := &Store{}
	s.Open()

	id, _ := s.Matchs().Create(protocol.Match{})

	itr, err := s.Matchs().CreateTransactionsChangeIterator(id)
	require.NoError(t, err)

	var (
		wg sync.WaitGroup

		givenTransactions = []protocol.Transaction{
			protocol.Transaction{Messages: []protocol.Message{}},
			protocol.Transaction{Messages: []protocol.Message{protocol.Message{}}},
		}
		receivedTrxChanges = []store.TransactionChange{}
	)

	wg.Add(4)

	go func() {
		var trxChange store.TransactionChange

		for itr.Next(&trxChange) {
			receivedTrxChanges = append(receivedTrxChanges, trxChange)

			wg.Done()
		}
	}()

	for _, trx := range givenTransactions {
		require.NoError(t, s.Matchs().CreateTransaction(id, trx))
	}

	require.NoError(t, s.Matchs().UpdateCurrentTransactionHolder(id, "foo"))

	require.NoError(t, s.Matchs().AddMessageToCurrentTransaction(id, protocol.Message{}))

	wg.Wait()

	require.Len(t, receivedTrxChanges, 4)

	require.Empty(t, receivedTrxChanges[0].Old)
	require.Empty(t, receivedTrxChanges[0].New.Messages)

	require.Empty(t, receivedTrxChanges[1].Old)
	require.Len(t, receivedTrxChanges[1].New.Messages, 1)

	require.Len(t, receivedTrxChanges[2].Old.Messages, 1)
	require.Len(t, receivedTrxChanges[2].New.Messages, 1)
	require.Equal(t, receivedTrxChanges[2].New.Holder.Name, "foo")

	require.Len(t, receivedTrxChanges[3].Old.Messages, 1)
	require.Len(t, receivedTrxChanges[3].New.Messages, 2)

	m, _ := s.Matchs().Read(id)

	require.Len(t, m.Transactions, 2)
	require.Len(t, m.Transactions[len(m.Transactions)-1].Messages, 2)

	require.NoError(t, itr.Close())
}
