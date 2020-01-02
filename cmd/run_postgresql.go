// +build run_postgresql

package cmd

import (
	"plateau/store"
	"plateau/store/postgresql"
)

func newStore() store.Store {
	return &postgresql.Store{}
}
