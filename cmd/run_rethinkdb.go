// +build run_rethinkdb

package cmd

import (
	"plateau/store/rethinkdb"
)

func newStore() *rethinkdb.Store {
	return &rethinkdb.Store{}
}
