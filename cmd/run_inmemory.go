// +build run_inmemory

package cmd

import (
	"plateau/store/inmemory"
)

func newStore() *inmemory.Store {
	return &inmemory.Store{}
}
