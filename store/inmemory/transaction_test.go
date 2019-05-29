package inmemory

import (
	"plateau/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransactionConversion(t *testing.T) {
	t.Parallel()

	var (
		trx = &transaction{}

		players = []*protocol.Player{}
	)

	require.IsType(t, &protocol.Transaction{}, trx.toProtocolStruct(players))
	require.IsType(t, trx, transactionFromProtocolStruct(trx.toProtocolStruct(players)))
}
