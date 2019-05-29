package inmemory

import (
	"plateau/store"

	"github.com/gorilla/sessions"
	"github.com/spf13/cobra"
)

var sessionKey string

// Store ...
type Store struct {
	inMemory *inMemory

	playerStore  *playerStore
	matchStore   *matchStore
	sessionStore sessions.Store
}

// Open implements `store.Store` interface.
func (s *Store) Open() error {
	inm := &inMemory{}

	s.playerStore = &playerStore{inm}
	s.matchStore = newmatchStore(inm)
	s.sessionStore = sessions.NewCookieStore([]byte([]byte(sessionKey)))

	return nil
}

// Close implements `store.Store` interface.
func (s *Store) Close() error {
	return s.matchStore.close()
}

// RunCommandSetter implements `store.Store` interface.
func (s *Store) RunCommandSetter(runCmd *cobra.Command) {
	runCmd.
		Flags().
		StringVarP(&sessionKey, "session-key", "", sessionKey, `Session ("secret") key`)
	runCmd.MarkFlagRequired("session-key")
	// TODO: Add a switch to configure the session expiration (MaxAge).
}

// Players implements `store.Store` interface.
func (s *Store) Players() store.PlayerStore {
	return s.playerStore
}

// Matchs implements `store.Store` interface.
func (s *Store) Matchs() store.MatchStore {
	return s.matchStore
}

// Sessions implements `store.Store` interface.
func (s *Store) Sessions() sessions.Store {
	return s.sessionStore
}
