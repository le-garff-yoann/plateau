package store

import (
	"github.com/gorilla/sessions"
	"github.com/spf13/cobra"
)

// Store represents a specialized means of access to a datastore
// (a SQL database for example) with the method `Open()` and `Close()`.
//	- `RunCommandSetter()` is there to allow the configuration of the awning via new Cobra flags.
//	- `BeginTransaction()` returns a new `Transaction`.
//	- `CreateMatchNotificationsIterator()` returns a new `MatchNotificationsIterator`.
type Store interface {
	Open() error
	Close() error

	RunCommandSetter(*cobra.Command)

	Sessions() sessions.Store

	BeginTransaction(...TransactionScope) Transaction
	CreateMatchNotificationsIterator(id string) (MatchNotificationsIterator, error)
}
