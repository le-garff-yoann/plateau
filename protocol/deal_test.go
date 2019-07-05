package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDealFind(t *testing.T) {
	t.Parallel()

	deal := Deal{Messages: []Message{Message{Code: MDealCompleted}}}

	require.NotNil(t, deal.Find(func(m Message) bool {
		return m.Code == MDealCompleted
	}))
	require.Nil(t, deal.Find(func(m Message) bool {
		return m.Code == MDealAborded
	}))
}

func TestDealFindAll(t *testing.T) {
	t.Parallel()

	deal := Deal{Messages: []Message{
		Message{Code: MPlayerAccepts},
		Message{Code: MPlayerAccepts},
		Message{Code: MPlayerRefuses},
	}}

	require.Len(t, deal.FindAll(func(m Message) bool {
		return m.Code == MPlayerAccepts
	}), 2)
	require.Len(t, deal.FindAll(func(m Message) bool {
		return m.Code == MPlayerRefuses
	}), 1)
	require.Empty(t, deal.FindAll(func(m Message) bool {
		return m.Code == MDealCompleted
	}))
}

func TestDealFindByMessageCode(t *testing.T) {
	t.Parallel()

	deal := Deal{Messages: []Message{Message{Code: MDealCompleted}}}

	require.NotNil(t, deal.FindByMessageCode(MDealCompleted))
	require.Nil(t, deal.FindByMessageCode(MDealAborded))
}

func TestDealFindAllByMessageCode(t *testing.T) {
	t.Parallel()

	deal := Deal{Messages: []Message{
		Message{Code: MPlayerAccepts},
		Message{Code: MPlayerAccepts},
		Message{Code: MPlayerRefuses},
	}}

	require.Len(t, deal.FindAllByMessageCode(MPlayerAccepts), 2)
	require.Len(t, deal.FindAllByMessageCode(MPlayerRefuses), 1)
	require.Empty(t, deal.FindAllByMessageCode(MDealCompleted))
}

func TestDealIsActive(t *testing.T) {
	t.Parallel()

	require.True(t, (&Deal{Messages: []Message{}}).IsActive())
	require.False(t, (&Deal{Messages: []Message{Message{Code: MDealCompleted}}}).IsActive())
	require.False(t, (&Deal{Messages: []Message{Message{Code: MDealAborded}}}).IsActive())
}

func TestWithMessagesConcealed(t *testing.T) {
	t.Parallel()

	deal := Deal{Messages: []Message{Message{
		Code: MPlayerAccepts,
		Payload: ConcealedMessagePayload{
			AllowedNamesCode: []string{"foo"},
		},
	}}}

	concealedDeal := deal.WithMessagesConcealed("bar")

	require.Len(t, deal.Messages, len(concealedDeal.Messages))
	require.NotEqual(t, deal.Messages[0].Code, concealedDeal.Messages[0].Code)
}
