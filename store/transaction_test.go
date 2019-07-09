package store

import (
	"plateau/protocol"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransactionScope(t *testing.T) {
	t.Parallel()

	trnScope := &TransactionScope{
		Mode:    TSReadMode,
		Subject: protocol.Player{Name: "foo"},
	}

	require.Equal(t, "foo", trnScope.IsSubjectPlayer().Name)
	require.Nil(t, trnScope.IsSubjectMatch())
	require.False(t, trnScope.IsSubjectAll())

	trnScope.Subject = protocol.Match{ID: "foo"}
	require.Equal(t, "foo", trnScope.IsSubjectMatch().ID)
	require.Nil(t, trnScope.IsSubjectPlayer())
	require.False(t, trnScope.IsSubjectAll())

	trnScope.Subject = nil
	require.True(t, trnScope.IsSubjectAll())
}
