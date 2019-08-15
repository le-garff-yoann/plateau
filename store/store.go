package store

import (
	"github.com/spf13/cobra"
)

// Store represents a specialized means of access to a datastore
// (a SQL database for example) with the method `Open()` and `Close()`.
//	- `RunCommandSetter()` is there to allow the configuration of the awning via new Cobra flags.
//	- `BeginTransaction()` returns a new `Transaction`.
//	- `RegisterNotificationsChannel()` registers a channel to receive store notifications.
//	- `UnregisterNotificationsChannel()` unregisters a channel registered with `RegisterNotificationsChannel()`.
type Store interface {
	Open() error
	Close() error

	RunCommandSetter(*cobra.Command)

	BeginTransaction(...TransactionScope) (Transaction, error)

	RegisterNotificationsChannel(chan interface{}) error
	UnregisterNotificationsChannel(chan interface{}) error
}
