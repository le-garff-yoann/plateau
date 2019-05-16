package store

import (
	"plateau/event"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEventContainerIsLegal(t *testing.T) {
	t.Parallel()

	ec := &EventContainer{}
	require.False(t, ec.IsLegal())

	anonymousEvents := []event.Event{
		event.EIllegal, event.EInternalError,
		event.EGameWantToStart, event.EGameStarts, event.EGameEnds,
		event.EPlayerConnects, event.EPlayerDisconnects, event.EPlayerReconnects,
		event.EPlayerJoins, event.EPlayerLeaves, event.EPlayerSurrenders,
	}

	for _, e := range anonymousEvents {
		require.True(t, (&EventContainer{Event: e}).IsLegal())
		require.False(t, (&EventContainer{Event: e, Emitter: &Player{}}).IsLegal())
	}

	userEmitterEvent := []event.Event{
		event.EListEvents, event.EPlayerWantToJoin, event.EPlayerWantToLeave, event.EPlayerWantToSurrender,
	}

	for _, e := range userEmitterEvent {
		require.True(t, (&EventContainer{Event: e}).IsLegal())
		require.True(t, (&EventContainer{Event: e, Emitter: &Player{}}).IsLegal())
	}
}

func TestEventContainerBroadcaster(t *testing.T) {
	t.Parallel()

	var (
		wg sync.WaitGroup

		br = NewEventContainerBroadcaster()

		i = 100
		a = 0
		b = 0

		recvA, _ = br.Recv()
		recvB, _ = br.Recv()
	)

	require.Equal(t, 2, br.CountReceivers())

	go br.Run()

	wg.Add(i * 2)
	go func() {
		for {
			select {
			case <-recvA:
				a++
			case <-recvB:
				b++
			}

			wg.Done()
		}
	}()

	for j := 0; j < i; j++ {
		br.Emitter <- EventContainer{}
	}

	wg.Wait()

	require.Equal(t, i, a)
	require.Equal(t, i, b)
}
