package store

import (
	"plateau/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchDealChange(t *testing.T) {
	t.Parallel()

	dealChange := DealChange{
		Old: &protocol.Deal{
			Messages: []protocol.Message{},
		},
		New: &protocol.Deal{
			Messages: []protocol.Message{protocol.Message{}},
		},
	}

	require.Len(t, dealChange.NewMessages(), 1)
}
