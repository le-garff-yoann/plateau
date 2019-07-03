// +build run_inmemory

package cmd

import (
	"plateau/store"
	"plateau/store/inmemory"
)

func newStore() store.Store {
	return &inmemory.Store{}
}
