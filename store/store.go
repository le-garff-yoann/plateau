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

	Sessions() sessions.Store

	BeginTransaction() Transaction
	CreateDealsChangeIterator(id string) (DealChangeIterator, error)
}
