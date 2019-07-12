package inmemory

import (
	"plateau/protocol"
	"plateau/store"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDealConversion(t *testing.T) {
	t.Parallel()

	var (
		deal = &deal{}

		players = []*protocol.Player{}
	)

	require.IsType(t, &protocol.Deal{}, deal.toProtocolStruct(players))
	require.IsType(t, deal, dealFromProtocolStruct(deal.toProtocolStruct(players)))
}

func TestCreateDealsChangeIterator(t *testing.T) {
	t.Parallel()

	s := &Store{}

	require.NoError(t, s.Open())
	defer func() {
		require.NoError(t, s.Close())
	}()

	trn := s.BeginTransaction()

	id, err := trn.MatchCreate(protocol.Match{})
	require.NoError(t, err)

	itr, err := s.CreateDealsChangeIterator(id)
	require.NoError(t, err)

	var (
		wg sync.WaitGroup

		givenDeals = []protocol.Deal{
			protocol.Deal{Messages: []protocol.Message{}},
			protocol.Deal{Messages: []protocol.Message{protocol.Message{}}},
		}
		receivedDealsChanges = []store.DealsChange{}
	)

	wg.Add(4)

	go func() {
		var dealChange store.DealsChange

		for itr.Next(&dealChange) {
			receivedDealsChanges = append(receivedDealsChanges, dealChange)

			wg.Done()
		}
	}()

	for _, deal := range givenDeals {
		require.NoError(t, trn.MatchCreateDeal(id, deal))
	}

	require.NoError(t, trn.MatchUpdateCurrentDealHolder(id, "foo"))

	require.NoError(t, trn.MatchAddMessageToCurrentDeal(id, protocol.Message{}))

	wg.Wait()

	require.Len(t, receivedDealsChanges, 4)

	require.Empty(t, receivedDealsChanges[0].Old)
	require.Empty(t, receivedDealsChanges[0].New.Messages)

	require.Empty(t, receivedDealsChanges[1].Old)
	require.Len(t, receivedDealsChanges[1].New.Messages, 1)

	require.Len(t, receivedDealsChanges[2].Old.Messages, 1)
	require.Len(t, receivedDealsChanges[2].New.Messages, 1)
	require.Equal(t, receivedDealsChanges[2].New.Holder.Name, "foo")

	require.Len(t, receivedDealsChanges[3].Old.Messages, 1)
	require.Len(t, receivedDealsChanges[3].New.Messages, 2)

	m, err := trn.MatchRead(id)
	require.NoError(t, err)

	require.Len(t, m.Deals, 2)
	require.Len(t, protocol.IndexDeals(m.Deals, 0).Messages, 2)

	require.NoError(t, itr.Close())
}
