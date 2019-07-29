package broadcaster

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBroadcaster(t *testing.T) {
	t.Parallel()

	var (
		wg sync.WaitGroup

		br = New()

		i = 10000
		a = 0
		b = 0

		chA, uuidA = br.Subscribe()
	)

	require.Len(t, br.subscribers, 1)

	go br.Run()

	chB, uuidB := br.Subscribe()
	require.Len(t, br.subscribers, 2)

	wg.Add(i * 2)
	go func() {
		for {
			select {
			case _, ok := <-chA:
				if !ok {
					return
				}

				a++
			case _, ok := <-chB:
				if !ok {
					return
				}

				b++
			}

			wg.Done()
		}
	}()

	for j := 0; j < i; j++ {
		go br.Submit(j)

		_, UUID := br.Subscribe()
		br.Unsubscribe(UUID)
	}

	wg.Wait()

	require.Equal(t, i, a)
	require.Equal(t, i, b)

	require.True(t, br.Unsubscribe(uuidA))
	require.True(t, br.Unsubscribe(uuidB))

	require.Empty(t, br.subscribers)
}
