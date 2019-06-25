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
	s.Open()

	trn := s.BeginTransaction()

	id, _ := trn.MatchCreate(protocol.Match{})

	itr, err := s.CreateDealsChangeIterator(id)
	require.NoError(t, err)

	var (
		wg sync.WaitGroup

		givenDeals = []protocol.Deal{
			protocol.Deal{Messages: []protocol.Message{}},
			protocol.Deal{Messages: []protocol.Message{protocol.Message{}}},
		}
		receivedDealChanges = []store.DealChange{}
	)

	wg.Add(4)

	go func() {
		var dealChange store.DealChange

		for itr.Next(&dealChange) {
			receivedDealChanges = append(receivedDealChanges, dealChange)

			wg.Done()
		}
	}()

	for _, deal := range givenDeals {
		require.NoError(t, trn.MatchCreateDeal(id, deal))
	}

	require.NoError(t, trn.MatchUpdateCurrentDealHolder(id, "foo"))

	require.NoError(t, trn.MatchAddMessageToCurrentDeal(id, protocol.Message{}))

	wg.Wait()

	require.Len(t, receivedDealChanges, 4)

	require.Empty(t, receivedDealChanges[0].Old)
	require.Empty(t, receivedDealChanges[0].New.Messages)

	require.Empty(t, receivedDealChanges[1].Old)
	require.Len(t, receivedDealChanges[1].New.Messages, 1)

	require.Len(t, receivedDealChanges[2].Old.Messages, 1)
	require.Len(t, receivedDealChanges[2].New.Messages, 1)
	require.Equal(t, receivedDealChanges[2].New.Holder.Name, "foo")

	require.Len(t, receivedDealChanges[3].Old.Messages, 1)
	require.Len(t, receivedDealChanges[3].New.Messages, 2)

	m, _ := trn.MatchRead(id)

	require.Len(t, m.Deals, 2)
	require.Len(t, m.Deals[len(m.Deals)-1].Messages, 2)

	require.NoError(t, itr.Close())
}
