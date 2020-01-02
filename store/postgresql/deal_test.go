package postgresql

import (
	"plateau/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDealConversion(t *testing.T) {
	t.Parallel()

	deal := &Deal{}

	require.IsType(t, &protocol.Deal{}, deal.toProtocolStruct())
	require.IsType(t, deal, dealFromProtocolStruct(deal.toProtocolStruct()))
}
