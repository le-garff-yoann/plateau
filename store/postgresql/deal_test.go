package postgresql

import (
	"testing"

	"github.com/le-garff-yoann/plateau/protocol"
	"github.com/stretchr/testify/require"
)

func TestDealConversion(t *testing.T) {
	t.Parallel()

	deal := &Deal{}

	require.IsType(t, &protocol.Deal{}, deal.toProtocolStruct())
	require.IsType(t, deal, dealFromProtocolStruct(deal.toProtocolStruct()))
}
