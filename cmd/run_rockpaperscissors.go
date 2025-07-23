//go:build run_rockpaperscissors
// +build run_rockpaperscissors

package cmd

import (
	"github.com/le-garff-yoann/plateau/game/rockpaperscissors"
	"github.com/le-garff-yoann/plateau/server"
)

func newGame() server.Game {
	return &rockpaperscissors.Game{}
}
