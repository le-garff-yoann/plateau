package store

import (
	"github.com/gorilla/sessions"
	"github.com/spf13/cobra"
)

// Store ...
type Store interface {
	Open() error
	Close() error

	RunCommandSetter(*cobra.Command)

	Players() PlayerStore
	Matchs() MatchStore
	Sessions() sessions.Store
}
