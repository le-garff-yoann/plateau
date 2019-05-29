package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransactionFind(t *testing.T) {
	t.Parallel()

	trx := Transaction{Messages: []Message{Message{MessageCode: MTransactionCompleted}}}

	require.NotNil(t, trx.Find(func(m Message) bool {
		return m.MessageCode == MTransactionCompleted
	}))
	require.Nil(t, trx.Find(func(m Message) bool {
		return m.MessageCode == MTransactionAborded
	}))
}

func TestTransactionFindAll(t *testing.T) {
	t.Parallel()

	trx := Transaction{Messages: []Message{
		Message{MessageCode: MPlayerAccepts},
		Message{MessageCode: MPlayerAccepts},
		Message{MessageCode: MPlayerRefuses},
	}}

	require.Len(t, trx.FindAll(func(m Message) bool {
		return m.MessageCode == MPlayerAccepts
	}), 2)
	require.Len(t, trx.FindAll(func(m Message) bool {
		return m.MessageCode == MPlayerRefuses
	}), 1)
	require.Empty(t, trx.FindAll(func(m Message) bool {
		return m.MessageCode == MTransactionCompleted
	}))
}

func TestTransactionFindByMessageCode(t *testing.T) {
	t.Parallel()

	trx := Transaction{Messages: []Message{Message{MessageCode: MTransactionCompleted}}}

	require.NotNil(t, trx.FindByMessageCode(MTransactionCompleted))
	require.Nil(t, trx.FindByMessageCode(MTransactionAborded))
}

func TestTransactionFindAllByMessageCode(t *testing.T) {
	t.Parallel()

	trx := Transaction{Messages: []Message{
		Message{MessageCode: MPlayerAccepts},
		Message{MessageCode: MPlayerAccepts},
		Message{MessageCode: MPlayerRefuses},
	}}

	require.Len(t, trx.FindAllByMessageCode(MPlayerAccepts), 2)
	require.Len(t, trx.FindAllByMessageCode(MPlayerRefuses), 1)
	require.Empty(t, trx.FindAllByMessageCode(MTransactionCompleted))
}

func TestTransactionIsActive(t *testing.T) {
	t.Parallel()

	require.True(t, (&Transaction{Messages: []Message{}}).IsActive())
	require.False(t, (&Transaction{Messages: []Message{Message{MessageCode: MTransactionCompleted}}}).IsActive())
	require.False(t, (&Transaction{Messages: []Message{Message{MessageCode: MTransactionAborded}}}).IsActive())
}
