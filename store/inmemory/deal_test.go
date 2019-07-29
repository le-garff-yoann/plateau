package inmemory

import (
	"plateau/protocol"
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
