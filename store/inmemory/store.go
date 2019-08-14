package inmemory

import (
	"plateau/broadcaster"
	"plateau/store"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/spf13/cobra"
)

var (
	sessionKeys = []string{""}

	sessionMaxAge = 86400 * 30
)

// Store implements the `store.Store` interface.
type Store struct {
	inMemoryMux sync.Mutex

	inMemory                      *inMemory
	matchNotificationsBroadcaster *broadcaster.Broadcaster

	sessionStore sessions.Store
}

// Open implements the `store.Store` interface.
func (s *Store) Open() error {
	s.inMemory = &inMemory{}

	s.matchNotificationsBroadcaster = broadcaster.New()
	go s.matchNotificationsBroadcaster.Run()

	var byteSessionKeys [][]byte
	for _, sessionKey := range sessionKeys {
		byteSessionKeys = append(byteSessionKeys, []byte(sessionKey))
	}

	sessionStore := sessions.NewCookieStore(byteSessionKeys...)
	sessionStore.MaxAge(sessionMaxAge)

	s.sessionStore = sessionStore

	return nil
}

// Close implements the `store.Store` interface.
func (s *Store) Close() error {
	s.matchNotificationsBroadcaster.Done()

	return nil
}

// RunCommandSetter implements the `store.Store` interface.
func (s *Store) RunCommandSetter(runCmd *cobra.Command) {
	runCmd.
		Flags().
		StringArrayVar(&sessionKeys, "session-key", sessionKeys, `Session ("secret") key`)
	runCmd.MarkFlagRequired("session-key")

	runCmd.
		Flags().
		IntVar(&sessionMaxAge, "session-max-age", sessionMaxAge, "Sets the maximum duration of cookies in seconds")
}

// Sessions implements the `store.Store` interface.
func (s *Store) Sessions() sessions.Store {
	return s.sessionStore
}

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
