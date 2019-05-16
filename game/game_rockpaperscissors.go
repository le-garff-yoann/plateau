// +build game_rockpaperscissors

package game

import (
	"plateau/game/rockpaperscissors"
)

// New ...
func New() Game {
	return &rockpaperscissors.Game{}
}
