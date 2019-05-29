package inmemory

import (
	"plateau/protocol"
	"sync"
)

type inMemory struct {
	mux sync.RWMutex

	players []*protocol.Player
	matchs  []*match
}

func (s *inMemory) player(name string) (player *protocol.Player) {
	for _, p := range s.players {
		if p.Name == name {
			player = p

			break
		}
	}

	return player
}

func (s *inMemory) match(id string) (match *match) {
	for _, m := range s.matchs {
		if m.ID == id {
			match = m

			break
		}
	}

	return match
}
