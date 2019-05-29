package store

import (
	"plateau/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchTransactionChange(t *testing.T) {
	t.Parallel()

	trxChange := TransactionChange{
		Old: &protocol.Transaction{
			Messages: []protocol.Message{},
		},
		New: &protocol.Transaction{
			Messages: []protocol.Message{protocol.Message{}},
		},
	}

	require.Len(t, trxChange.NewMessages(), 1)
}
