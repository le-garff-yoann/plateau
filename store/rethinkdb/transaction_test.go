package rethinkdb

import (
	"plateau/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransactionConversion(t *testing.T) {
	t.Parallel()

	trx := &transaction{}
	require.IsType(t, &protocol.Transaction{}, trx.toProtocolStruct())
	require.IsType(t, trx, transactionFromProtocolStruct(trx.toProtocolStruct()))
}
