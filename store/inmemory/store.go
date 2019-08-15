package inmemory

import (
	"plateau/broadcaster"
	"plateau/store"
	"sync"

	"github.com/spf13/cobra"
)

var (
	sessionKeys   = []string{""}
	sessionMaxAge = 86400 * 30
)

// Store implements the `store.Store` interface.
type Store struct {
	inMemoryMux sync.Mutex

	inMemory                      *inMemory
	matchNotificationsBroadcaster *broadcaster.Broadcaster
}

// Open implements the `store.Store` interface.
func (s *Store) Open() error {
	s.inMemory = &inMemory{}

	s.matchNotificationsBroadcaster = broadcaster.New()
	go s.matchNotificationsBroadcaster.Run()

	return nil
}

// Close implements the `store.Store` interface.
func (s *Store) Close() error {
	s.matchNotificationsBroadcaster.Done()

	return nil
}

// RunCommandSetter implements the `store.Store` interface.
func (s *Store) RunCommandSetter(runCmd *cobra.Command) {}

// BeginTransaction implements the `store.Store` interface.
func (s *Store) BeginTransaction(scopes ...store.TransactionScope) store.Transaction {
	s.inMemoryMux.Lock()

	return &Transaction{
		inMemory:     s.inMemory,
		inMemoryCopy: s.inMemory.Copy(),
		closed:       false,
		commitCb: func(trn *Transaction) {
			for _, n := range trn.matchNotifications {
				s.matchNotificationsBroadcaster.Submit(n)
			}
		},
		done: func(_ *Transaction) { s.inMemoryMux.Unlock() },
	}
}

// RegisterNotificationsChannel implements the `store.Store` interface.
func (s *Store) RegisterNotificationsChannel(ch chan interface{}) error {
	s.matchNotificationsBroadcaster.Register(ch)

	return nil
}

// UnregisterNotificationsChannel implements the `store.Store` interface.
func (s *Store) UnregisterNotificationsChannel(ch chan interface{}) error {
	s.matchNotificationsBroadcaster.Unregister(ch)

	return nil
}
