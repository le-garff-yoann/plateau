package inmemory

import (
	"plateau/protocol"

	"github.com/ulule/deepcopier"
)

type inMemory struct {
	Players []*protocol.Player
	Matchs  []*match
}

func (s *inMemory) Player(name string) (player *protocol.Player) {
	for _, p := range s.Players {
		if p.Name == name {
			player = p

			break
		}
	}

	return player
}

func (s *inMemory) Match(id string) (match *match) {
	for _, m := range s.Matchs {
		if m.ID == id {
			match = m

			break
		}
	}

	return match
}

func (s *inMemory) Copy() *inMemory {
	var inMemoryCopy inMemory
	deepcopier.Copy(s).To(&inMemoryCopy)

	return &inMemoryCopy
}
