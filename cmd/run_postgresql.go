//go:build run_postgresql
// +build run_postgresql

package cmd

import (
	"github.com/le-garff-yoann/plateau/store"
	"github.com/le-garff-yoann/plateau/store/postgresql"
)

func newStore() store.Store {
	return &postgresql.Store{}
}
