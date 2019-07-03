// +build run_rockpaperscissors

package cmd

import (
	"plateau/game/rockpaperscissors"
	"plateau/server"
)

func newGame() server.Game {
	return &rockpaperscissors.Game{}
}
