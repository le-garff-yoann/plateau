package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotificationContainerString(t *testing.T) {
	t.Parallel()

	require.NotPanics(t, func() {
		t.Log((&NotificationContainer{
			Notification: NDealChange,
			Body:         []string{"foo", "bar"},
		}).String())
	})
}
