package server

import (
	"errors"
	"plateau/broadcaster"
	"plateau/protocol"
	"plateau/server/response/body"
	"plateau/store"
)

func (s *Server) guardRuntime(matchID string) (*matchRuntime, error) {
	s.matchRuntimesMux.Lock()
	defer s.matchRuntimesMux.Unlock()

	if _, ok := s.matchRuntimes[matchID]; !ok {
		iterator, err := s.store.CreateMatchNotificationsIterator(matchID)
		if err != nil {
			return nil, err
		}

		mRuntime := &matchRuntime{
			game:                          s.game,
			matchID:                       matchID,
			matchNotificationsBroadcaster: broadcaster.New(),
			done:                          make(chan int),
		}

		s.matchRuntimes[matchID] = mRuntime

		go mRuntime.matchNotificationsBroadcaster.Run()

		go func() {
			var matchNotification store.MatchNotification

			for iterator.Next(&matchNotification) {
				func(r *matchRuntime) {
					r.matchNotificationsBroadcaster.Submit(matchNotification)
				}(mRuntime)
			}
		}()

		go func(r *matchRuntime) {
			<-r.done

			r.matchNotificationsBroadcaster.Done()

			if iterator.Close() == nil {
				return
			}
		}(mRuntime)
	}

	s.matchRuntimes[matchID].guard++

	return s.matchRuntimes[matchID], nil
}

func (s *Server) unguardRuntime(matchID string) bool {
	s.matchRuntimesMux.Lock()
	defer s.matchRuntimesMux.Unlock()

	r, ok := s.matchRuntimes[matchID]
	if ok {
		r.guard--
		if r.guard == 0 {
			s.matchRuntimes[matchID].done <- 0

			delete(s.matchRuntimes, matchID)
		}
	}

	return ok
}

type matchRuntime struct {
	game Game

	matchID string
	guard   int

	matchNotificationsBroadcaster *broadcaster.Broadcaster

	done chan int
}

func (s *matchRuntime) reqContainerHandler(trn store.Transaction, reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
	match, err := trn.MatchRead(s.matchID)
	if err != nil {
		trn.Abort()

		return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
	}

	reqContainer.Match = match

	ctx := baseContext(trn, reqContainer)

	if !reqContainer.Match.IsEnded() {
		for _, deal := range reqContainer.Match.Deals {
			if deal.FindByMessageCode(protocol.MPlayerWantToStartTheGame) != nil && deal.FindByMessageCode(protocol.MDealCompleted) != nil {
				return baseContext(trn, reqContainer).Complete(
					s.game.Context(trn, reqContainer),
				).handle(reqContainer)
			}
		}

		currentDeal := protocol.IndexDeals(reqContainer.Match.Deals, 0)

		if currentDeal == nil || !currentDeal.IsActive() {
			func() {
				for _, p := range reqContainer.Match.Players {
					if p.Name == reqContainer.Player.Name {
						ctx.Complete(leaveContext(trn, reqContainer)).Complete(wantToStartMatchContext(trn, reqContainer))

						return
					}
				}

				ctx.Complete(joinContext(trn, reqContainer))
			}()
		} else {
			if currentDeal.Holder.Name == reqContainer.Player.Name {
				ctx.Complete(askToStartMatchContext(trn, reqContainer))
			}
		}
	}

	return ctx.handle(reqContainer)
}

func baseContext(trn store.Transaction, reqContainer *protocol.RequestContainer) *Context {
	ctx := NewContext()
	ctx.
		OnNotImplemented(func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			trn.Commit()

			return &protocol.ResponseContainer{Response: protocol.ResNotImplemented}
		}).
		On(protocol.ReqListRequests, func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			return &protocol.ResponseContainer{
				Response: protocol.ResOK,
				Body:     ctx.Requests(),
			}
		}).
		After(func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			trn.Commit()

			return nil
		})

	return ctx
}

func joinContext(trn store.Transaction, reqContainer *protocol.RequestContainer) *Context {
	return NewContext().
		On(protocol.ReqPlayerWantToJoin, func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			if err := trn.MatchPlayerJoins(reqContainer.Match.ID, reqContainer.Player.Name); err != nil {
				defer trn.Abort()

				if _, ok := err.(store.PlayerParticipationError); ok {
					return &protocol.ResponseContainer{Response: protocol.ResForbidden, Body: body.New().Ko(err)}
				}

				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		})
}

func leaveContext(trn store.Transaction, reqContainer *protocol.RequestContainer) *Context {
	return NewContext().
		On(protocol.ReqPlayerWantToLeave, func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			if err := trn.MatchPlayerLeaves(reqContainer.Match.ID, reqContainer.Player.Name); err != nil {
				defer trn.Abort()

				if _, ok := err.(store.PlayerParticipationError); ok {
					return &protocol.ResponseContainer{Response: protocol.ResForbidden, Body: body.New().Ko(err)}
				}

				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		})
}

func wantToStartMatchContext(trn store.Transaction, reqContainer *protocol.RequestContainer) *Context {
	return NewContext().
		On(protocol.ReqPlayerWantToStartTheMatch, func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			if !reqContainer.Match.IsFull() {
				return &protocol.ResponseContainer{
					Response: protocol.ResForbidden,
					Body:     body.New().Ko(errors.New("There are not enough players")),
				}
			}

			if err := trn.MatchCreateDeal(reqContainer.Match.ID, protocol.Deal{Holder: *reqContainer.Player, Messages: []protocol.Message{protocol.Message{Code: protocol.MPlayerWantToStartTheGame}}}); err != nil {
				trn.Abort()

				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		})
}

func askToStartMatchContext(trn store.Transaction, reqContainer *protocol.RequestContainer) *Context {
	return NewContext().
		On(protocol.ReqPlayerAccepts, func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			if err := trn.MatchAddMessageToCurrentDeal(reqContainer.Match.ID, protocol.Message{Code: protocol.MPlayerAccepts, Payload: reqContainer.Player.Name}); err != nil {
				trn.Abort()

				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			match, err := trn.MatchRead(reqContainer.Match.ID)
			if err != nil {
				trn.Abort()

				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			var OKPlayersName []string
			for _, m := range protocol.IndexDeals(match.Deals, 0).FindAllByMessageCode(protocol.MPlayerAccepts) {
				OKPlayersName = append(OKPlayersName, m.Payload.(string))
			}

			for _, player := range match.Players {
				if !func() bool {
					for _, OKPlayerName := range OKPlayersName {
						if player.Name == OKPlayerName {
							return true
						}
					}

					return false
				}() {
					if err := trn.MatchUpdateCurrentDealHolder(reqContainer.Match.ID, player.Name); err != nil {
						trn.Abort()

						return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
					}

					return &protocol.ResponseContainer{Response: protocol.ResOK}
				}
			}

			if err := trn.MatchAddMessageToCurrentDeal(reqContainer.Match.ID, protocol.Message{Code: protocol.MDealCompleted}); err != nil {
				trn.Abort()

				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		}).
		On(protocol.ReqPlayerRefuses, func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			if err := trn.MatchAddMessageToCurrentDeal(reqContainer.Match.ID, protocol.Message{Code: protocol.MDealAborded}); err != nil {
				trn.Abort()

				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		})
}
