package server

import (
	"log"
	"net/http"
	"plateau/event"
	"plateau/game"
	"plateau/store"

	"github.com/gorilla/mux"
)

func (s *Server) connectMatchHandler(w http.ResponseWriter, r *http.Request) {
	session, err := s.store.Sessions().Get(r, ServerName)
	if err != nil {
		return
	}

	v := mux.Vars(r)

	c, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer func() {
		c.Close()
	}()

	gameEngine := game.New()

	if err := gameEngine.Init(); err != nil {
		if c.WriteJSON(store.EventContainer{Event: event.EInternalError}) != nil {
			return
		}
	}

	matchID := v["id"]

	player, err := s.store.Players().Read(session.Values["username"].(string))
	if err != nil {
		if c.WriteJSON(store.EventContainer{Event: event.EInternalError}) != nil {
			return
		}
	}

	done := make(chan int)

	recv, brUUID, err := s.recvEventContainerBroadcaster(matchID)
	if err != nil {
		if c.WriteJSON(store.EventContainer{Event: event.EInternalError}) != nil {
			return
		}
	}
	defer func() {
		s.removeRecvEventContainerBroadcaster(matchID, *brUUID)

		s.store.Matchs().CreateEventContainer(matchID, store.EventContainer{Event: event.EPlayerDisconnects})
	}()

	go func() {
		for {
			select {
			case ec := <-recv:
				if c.WriteJSON(store.EventContainer{Event: event.EProcessed}) != nil {
					return
				}

				if c.WriteJSON(ec) != nil {
					continue
				}
			case <-done:
				return
			}
		}
	}()

	defer func() {
		done <- 0
	}()

	for {
		var eventContainer store.EventContainer

		if c.ReadJSON(&eventContainer) != nil {
			if c.WriteJSON(store.EventContainer{Event: event.EIllegal}) != nil {
				return
			}

			continue
		}

		eventContainer.Emitter = player

		if !eventContainer.IsLegal() {
			if c.WriteJSON(store.EventContainer{Event: event.EIllegal}) != nil {
				return
			}
		}

		log.Printf(eventContainer.String())

		if err := s.store.Matchs().CreateEventContainer(matchID, eventContainer); err != nil {
			if c.WriteJSON(store.EventContainer{Event: event.EIllegal}) != nil {
				return
			}

			continue
		}

		// var (
		// 	err error

		// 	isPlayerIn = false
		// )

		// match, err := s.store.Matchs().Read(matchID)
		// if err != nil {
		// 	if c.WriteJSON(store.EventContainer{Event: event.EInternalError}) != nil {
		// 		return
		// 	}

		// 	continue
		// }

		// for _, p := range match.Players {
		// 	if p.Name == eventContainer.Emitter.Name {
		// 		isPlayerIn = true

		// 		break
		// 	}
		// }

		// if g.OnEvent(match, &eventContainer) != nil {
		// 	if c.WriteJSON(store.EventContainer{Event: event.EIllegal}) != nil {
		// 		return
		// 	}

		// 	continue
		// }

		// switch eventContainer.Event {
		// case event.EPlayerWantToJoin:
		// 	if isPlayerIn {
		// 		if c.WriteJSON(store.EventContainer{Event: event.EIllegal}) != nil {
		// 			return
		// 		}

		// 		continue
		// 	}

		// 	if err := s.store.Matchs().AddPlayer(match.ID, eventContainer.Emitter.Name); err != nil {
		// 		if c.WriteJSON(store.EventContainer{Event: event.EInternalError}) != nil {
		// 			return
		// 		}

		// 		continue
		// 	}

		// 	eventContainer.Subjects = append(eventContainer.Subjects, eventContainer.Emitter)

		// 	pairedEventContainer := &eventContainer
		// 	pairedEventContainer.Event = event.EPlayerJoins

		// 	g.OnEvent(match, pairedEventContainer)
		// case event.EPlayerWantToLeave:
		// 	if !isPlayerIn || match.Running {
		// 		if c.WriteJSON(store.EventContainer{Event: event.EIllegal}) != nil {
		// 			return
		// 		}

		// 		continue
		// 	}

		// 	if s.store.Matchs().RemovePlayer(match.ID, eventContainer.Emitter.Name) != nil {
		// 		if c.WriteJSON(store.EventContainer{Event: event.EInternalError}) != nil {
		// 			return
		// 		}

		// 		continue
		// 	}

		// 	eventContainer.Subjects = append(eventContainer.Subjects, eventContainer.Emitter)

		// 	pairedEventContainer := &eventContainer
		// 	pairedEventContainer.Event = event.EPlayerLeaves

		// 	g.OnEvent(match, pairedEventContainer)
		// }

		// if c.WriteJSON(store.EventContainer{Event: event.EProcessed}) != nil {
		// 	return
		// }
	}
}
