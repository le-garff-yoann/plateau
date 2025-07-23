//go:build run_inmemory
// +build run_inmemory

package cmd

import (
	"github.com/le-garff-yoann/plateau/store"
	"github.com/le-garff-yoann/plateau/store/inmemory"
)

func newStore() store.Store {
	return &inmemory.Store{}
}
