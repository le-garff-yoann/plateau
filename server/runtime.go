package server

import (
	"errors"
	"plateau/broadcaster"
	"plateau/protocol"
	"plateau/server/response/body"
	"plateau/store"
	"sync"

	"github.com/thoas/go-funk"
)

func (s *Server) guardRuntime(matchID string) (*MatchRuntime, error) {
	if _, ok := s.matchRuntimes[matchID]; !ok {
		iterator, err := s.store.CreateDealsChangeIterator(matchID)
		if err != nil {
			return nil, err
		}

		s.matchmatchRuntimesMux.Lock()
		defer s.matchmatchRuntimesMux.Unlock()
		s.matchRuntimes[matchID] = &MatchRuntime{
			game:                    s.game,
			matchID:                 matchID,
			dealsChanges:            []store.DealChange{},
			dealsChangesBroadcaster: broadcaster.New(),
			done:                    make(chan int),
		}

		go s.matchRuntimes[matchID].dealsChangesBroadcaster.Run()

		go func() {
			var dealChange store.DealChange

			for iterator.Next(&dealChange) {
				func(r *MatchRuntime) {
					r.dealsChangesMux.Lock()
					defer r.dealsChangesMux.Unlock()

					r.dealsChanges = append(r.dealsChanges, dealChange)

					r.dealsChangesBroadcaster.Submit(protocol.NotificationContainer{
						Notification: protocol.NDealChange,
						Body:         dealChange,
					})
				}(s.matchRuntimes[matchID])
			}
		}()

		go func(r *MatchRuntime) {
			<-r.done

			r.dealsChangesBroadcaster.Done()

			if iterator.Close() == nil {
				return
			}
		}(s.matchRuntimes[matchID])
	}

	s.matchRuntimes[matchID].guard++

	return s.matchRuntimes[matchID], nil
}

func (s *Server) unguardRuntime(matchID string) bool {
	r, ok := s.matchRuntimes[matchID]

	if ok {
		s.matchmatchRuntimesMux.Lock()
		defer s.matchmatchRuntimesMux.Unlock()

		r.guard--
		if r.guard == 0 {
			s.matchRuntimes[matchID].done <- 0

			delete(s.matchRuntimes, matchID)
		}
	}

	return ok
}

// MatchRuntime ...
type MatchRuntime struct {
	game Game

	matchID string
	guard   int

	dealsChangesMux         sync.RWMutex
	dealsChanges            []store.DealChange
	dealsChangesBroadcaster *broadcaster.Broadcaster

	done chan int
}

// DealsChanges ...
func (s *MatchRuntime) DealsChanges() []store.DealChange {
	s.dealsChangesMux.RLock()
	defer s.dealsChangesMux.RUnlock()

	return s.dealsChanges
}

func (s *MatchRuntime) requestContainerHandler(trn store.Transaction, requestContainer *protocol.RequestContainer) *protocol.ResponseContainer {
	match, err := trn.MatchRead(s.matchID)
	if err != nil {
		trn.Abort()

		return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
	}

	requestContainer.Match = match

	ctx := baseContext(requestContainer)

	if requestContainer.Match.EndedAt == nil {
		for _, deal := range requestContainer.Match.Deals {
			if deal.FindByMessageCode(protocol.MPlayerWantToStartTheGame) != nil && deal.FindByMessageCode(protocol.MDealCompleted) != nil {
				return baseContext(requestContainer).Complete(
					s.game.Context(s, trn, requestContainer),
				).handle(s, requestContainer)
			}
		}

		currentDeal := protocol.IndexDeals(requestContainer.Match.Deals, 0)

		if currentDeal == nil || !currentDeal.IsActive() {
			if funk.Contains(requestContainer.Match.Players, *requestContainer.Player) {
				ctx.Complete(leaveContext(trn, requestContainer)).Complete(wantToStartMatchContext(trn, requestContainer))
			} else {
				ctx.Complete(joinContext(trn, requestContainer))
			}
		} else {
			if currentDeal.Holder.Name == requestContainer.Player.Name {
				ctx.Complete(askToStartMatchContext(trn, requestContainer))
			}
		}
	}

	trn.Commit()

	return ctx.handle(s, requestContainer)
}

func baseContext(requestContainer *protocol.RequestContainer) *Context {
	ctx := NewContext()
	ctx.
		On(protocol.ReqGetCurrentDeal, func(matchRuntime *MatchRuntime, requestContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			return &protocol.ResponseContainer{
				Response: protocol.ResOK,
				Body:     protocol.IndexDeals(requestContainer.Match.Deals, 0),
			}
		}).
		On(protocol.ReqListRequests, func(matchRuntime *MatchRuntime, requestContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			return &protocol.ResponseContainer{
				Response: protocol.ResOK,
				Body:     ctx.Requests(),
			}
		})

	return ctx
}

func joinContext(trn store.Transaction, requestContainer *protocol.RequestContainer) *Context {
	return NewContext().
		On(protocol.ReqPlayerWantToJoin, func(matchRuntime *MatchRuntime, requestContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			if err := trn.MatchPlayerJoins(matchRuntime.matchID, requestContainer.Player.Name); err != nil {
				defer trn.Abort()

				if _, ok := err.(store.PlayerParticipationError); ok {
					return &protocol.ResponseContainer{Response: protocol.ResForbidden, Body: body.New().Ko(err)}
				}

				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		})
}

func leaveContext(trn store.Transaction, requestContainer *protocol.RequestContainer) *Context {
	return NewContext().
		On(protocol.ReqPlayerWantToLeave, func(matchRuntime *MatchRuntime, requestContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			if err := trn.MatchPlayerLeaves(matchRuntime.matchID, requestContainer.Player.Name); err != nil {
				defer trn.Abort()

				if _, ok := err.(store.PlayerParticipationError); ok {
					return &protocol.ResponseContainer{Response: protocol.ResForbidden, Body: body.New().Ko(err)}
				}

				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		})
}

func wantToStartMatchContext(trn store.Transaction, requestContainer *protocol.RequestContainer) *Context {
	return NewContext().
		On(protocol.ReqPlayerWantToStartTheGame, func(matchRuntime *MatchRuntime, requestContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			if !requestContainer.Match.IsFull() {
				return &protocol.ResponseContainer{
					Response: protocol.ResForbidden,
					Body:     body.New().Ko(errors.New("There are not enough players")),
				}
			}

			if err := trn.MatchCreateDeal(matchRuntime.matchID, protocol.Deal{Holder: *requestContainer.Player, Messages: []protocol.Message{protocol.Message{MessageCode: protocol.MPlayerWantToStartTheGame}}}); err != nil {
				trn.Abort()

				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		})
}

func askToStartMatchContext(trn store.Transaction, requestContainer *protocol.RequestContainer) *Context {
	return NewContext().
		On(protocol.ReqPlayerAccepts, func(matchRuntime *MatchRuntime, requestContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			if err := trn.MatchAddMessageToCurrentDeal(matchRuntime.matchID, protocol.Message{MessageCode: protocol.MPlayerAccepts, Payload: requestContainer.Player.Name}); err != nil {
				trn.Abort()

				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			match, err := trn.MatchRead(matchRuntime.matchID)
			if err != nil {
				trn.Abort()

				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			var OKPlayersName []string
			for _, m := range protocol.IndexDeals(match.Deals, 0).FindAllByMessageCode(protocol.MPlayerAccepts) {
				OKPlayersName = append(OKPlayersName, m.Payload.(string))
			}

			for _, player := range match.Players {
				if !funk.Contains(OKPlayersName, player.Name) {
					if err := trn.MatchUpdateCurrentDealHolder(matchRuntime.matchID, player.Name); err != nil {
						trn.Abort()

						return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
					}

					return &protocol.ResponseContainer{Response: protocol.ResOK}
				}
			}

			if err := trn.MatchAddMessageToCurrentDeal(matchRuntime.matchID, protocol.Message{MessageCode: protocol.MDealCompleted}); err != nil {
				trn.Abort()

				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		}).
		On(protocol.ReqPlayerRefuses, func(matchRuntime *MatchRuntime, requestContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			if err := trn.MatchAddMessageToCurrentDeal(matchRuntime.matchID, protocol.Message{MessageCode: protocol.MDealAborded}); err != nil {
				trn.Abort()

				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		})
}
