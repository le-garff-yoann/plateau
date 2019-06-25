package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDealFind(t *testing.T) {
	t.Parallel()

	deal := Deal{Messages: []Message{Message{MessageCode: MDealCompleted}}}

	require.NotNil(t, deal.Find(func(m Message) bool {
		return m.MessageCode == MDealCompleted
	}))
	require.Nil(t, deal.Find(func(m Message) bool {
		return m.MessageCode == MDealAborded
	}))
}

func TestDealFindAll(t *testing.T) {
	t.Parallel()

	deal := Deal{Messages: []Message{
		Message{MessageCode: MPlayerAccepts},
		Message{MessageCode: MPlayerAccepts},
		Message{MessageCode: MPlayerRefuses},
	}}

	require.Len(t, deal.FindAll(func(m Message) bool {
		return m.MessageCode == MPlayerAccepts
	}), 2)
	require.Len(t, deal.FindAll(func(m Message) bool {
		return m.MessageCode == MPlayerRefuses
	}), 1)
	require.Empty(t, deal.FindAll(func(m Message) bool {
		return m.MessageCode == MDealCompleted
	}))
}

func TestDealFindByMessageCode(t *testing.T) {
	t.Parallel()

	deal := Deal{Messages: []Message{Message{MessageCode: MDealCompleted}}}

	require.NotNil(t, deal.FindByMessageCode(MDealCompleted))
	require.Nil(t, deal.FindByMessageCode(MDealAborded))
}

func TestDealFindAllByMessageCode(t *testing.T) {
	t.Parallel()

	deal := Deal{Messages: []Message{
		Message{MessageCode: MPlayerAccepts},
		Message{MessageCode: MPlayerAccepts},
		Message{MessageCode: MPlayerRefuses},
	}}

	require.Len(t, deal.FindAllByMessageCode(MPlayerAccepts), 2)
	require.Len(t, deal.FindAllByMessageCode(MPlayerRefuses), 1)
	require.Empty(t, deal.FindAllByMessageCode(MDealCompleted))
}

func TestDealIsActive(t *testing.T) {
	t.Parallel()

	require.True(t, (&Deal{Messages: []Message{}}).IsActive())
	require.False(t, (&Deal{Messages: []Message{Message{MessageCode: MDealCompleted}}}).IsActive())
	require.False(t, (&Deal{Messages: []Message{Message{MessageCode: MDealAborded}}}).IsActive())
}
