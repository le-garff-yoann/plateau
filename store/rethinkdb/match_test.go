package rethinkdb

import (
	"plateau/store"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	rethinkdb "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

func TestMatchStoreConversion(t *testing.T) {
	t.Parallel()

	m := Match{}
	require.IsType(t, &store.Match{}, m.toStoreStruct())
	require.IsType(t, m, *matchFromStoreStruct(*m.toStoreStruct()))
}

func TestMatchList(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		matchStore = &MatchStore{mock}
	)

	mock.On(matchStore.listTerm()).Return([]string{
		"foo", "bar",
	}, nil)

	matchs, err := matchStore.List()
	require.NoError(t, err)
	require.Len(t, matchs, 2)

	mock.AssertExpectations(t)
}

func TestMatchCreate(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		matchStore = &MatchStore{mock}

		givenMatch = store.Match{ID: "foo"}
	)

	mock.On(matchStore.createTerm(givenMatch)).Return(
		rethinkdb.WriteResponse{GeneratedKeys: []string{givenMatch.ID}},
		nil)

	id, err := matchStore.Create(givenMatch)
	require.NoError(t, err)
	require.Equal(t, givenMatch.ID, id)

	mock.AssertExpectations(t)
}

func TestMatchRead(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		matchStore = &MatchStore{mock}

		givenID       = "foo"
		givenMatch    = Match{ID: givenID}
		expectedMatch = *givenMatch.toStoreStruct()
	)

	mock.On(matchStore.readTerm(givenID)).Return([]interface{}{
		givenMatch,
	}, nil)

	match, err := matchStore.Read(givenID)
	require.NoError(t, err)
	require.IsType(t, &expectedMatch, match)
	require.Equal(t, expectedMatch.ID, match.ID)

	mock.AssertExpectations(t)
}

func TestMatchEndedAt(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		matchStore = &MatchStore{mock}

		givenMatch = Match{ID: "foo"}
		givenTime  = time.Now()
	)

	mock.On(matchStore.endedAtTerm(givenMatch.ID, &givenTime)).Return(rethinkdb.WriteResponse{},
		nil)

	err := matchStore.EndedAt(givenMatch.ID, &givenTime)
	require.NoError(t, err)
}

func TestMatchAddPlayer(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		matchStore = &MatchStore{mock}

		givenMatchID    = "foo"
		givenPlayerName = "bar"
	)

	mock.On(matchStore.addPlayerTerm(givenMatchID, givenPlayerName)).Return(rethinkdb.WriteResponse{},
		nil)

	err := matchStore.AddPlayer(givenMatchID, givenPlayerName)
	require.NoError(t, err)
}

func TestMatchRemovePlayer(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		matchStore = &MatchStore{mock}

		givenMatchID    = "foo"
		givenPlayerName = "bar"
	)

	mock.On(matchStore.removePlayerTerm(givenMatchID, givenPlayerName)).Return(rethinkdb.WriteResponse{},
		nil)

	err := matchStore.RemovePlayer(givenMatchID, givenPlayerName)
	require.NoError(t, err)
}

func TestMatchRunning(t *testing.T) {
	t.Parallel()

	var (
		mock = rethinkdb.NewMock()

		matchStore = &MatchStore{mock}

		givenMatch   = Match{ID: "foo"}
		givenRunning = true
	)

	mock.On(matchStore.runningTerm(givenMatch.ID, givenRunning)).Return(rethinkdb.WriteResponse{},
		nil)

	err := matchStore.Running(givenMatch.ID, givenRunning)
	require.NoError(t, err)
}

func TestMatchCreateEventContainer(t *testing.T) {
	var (
		mock = rethinkdb.NewMock()

		matchStore = &MatchStore{mock}

		givenMatchID        = "foo"
		givenEventContainer = store.EventContainer{}
	)

	mock.On(matchStore.createEventContainerTerm(givenMatchID, givenEventContainer)).Return(rethinkdb.WriteResponse{},
		nil)

	err := matchStore.CreateEventContainer(givenMatchID, givenEventContainer)
	require.NoError(t, err)
}

func TestMatchCreateEventContainerBroadcaster(t *testing.T) {
	var (
		mock = rethinkdb.NewMock()

		matchStore = &MatchStore{mock}

		givenMatchID             = "foo"
		givenMatchChangeResponse = matchChangeResponse{
			OldValue: Match{EventContainers: []EventContainer{}},
			NewValue: Match{EventContainers: []EventContainer{EventContainer{}, EventContainer{}}},
		}
	)

	mock.On(matchStore.eventContainerChangesTerm(givenMatchID)).Return(givenMatchChangeResponse, nil)

	br, err := matchStore.CreateEventContainerBroadcaster(givenMatchID)
	require.NoError(t, err)

	var (
		wg sync.WaitGroup

		recv, _ = br.Recv()

		i = 0
	)

	go br.Run()

	wg.Add(len(givenMatchChangeResponse.NewValue.EventContainers))
	go func() {
		for {
			<-recv

			i++
			wg.Done()
		}
	}()

	wg.Wait()

	require.Equal(t, len(givenMatchChangeResponse.NewValue.EventContainers), i)
}
