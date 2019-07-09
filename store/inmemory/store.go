package inmemory

import (
	"plateau/broadcaster"
	"plateau/store"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/spf13/cobra"
	"github.com/ulule/deepcopier"
)

var sessionKey string

// Store implements the `store.Store` interface.
type Store struct {
	inMemoryMux sync.Mutex

	inMemory               *inMemory
	dealChangesBroadcaster *broadcaster.Broadcaster

	sessionStore sessions.Store
}

// Open implements the `store.Store` interface.
func (s *Store) Open() error {
	s.inMemory = &inMemory{}

	s.dealChangesBroadcaster = broadcaster.New()
	go s.dealChangesBroadcaster.Run()

	s.sessionStore = sessions.NewCookieStore([]byte([]byte(sessionKey)))

	return nil
}

// Close implements the `store.Store` interface.
func (s *Store) Close() error {
	s.dealChangesBroadcaster.Done()

	return nil
}

// RunCommandSetter implements the `store.Store` interface.
func (s *Store) RunCommandSetter(runCmd *cobra.Command) {
	runCmd.
		Flags().
		StringVarP(&sessionKey, "session-key", "", sessionKey, `Session ("secret") key`)
	runCmd.MarkFlagRequired("session-key")
	// TODO: Add a switch to configure the session expiration (MaxAge).
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
		dealChangeSubmitter: func(dealChange *store.DealChange) {
			var dealChangeCopy store.DealChange
			deepcopier.Copy(dealChange).To(&dealChangeCopy)

			s.dealChangesBroadcaster.Submit(dealChangeCopy)
		},
		closed: false,
		done:   func() { s.inMemoryMux.Unlock() },
	}
}
