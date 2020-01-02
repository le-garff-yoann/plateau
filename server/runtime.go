package server

import (
	"errors"
	"plateau/protocol"
	"plateau/server/response/body"
	"plateau/store"
)

func (s *Server) guardRuntime(matchID string) (*matchRuntime, error) {
	s.matchRuntimesMux.Lock()
	defer s.matchRuntimesMux.Unlock()

	if _, ok := s.matchRuntimes[matchID]; !ok {
		mRuntime := &matchRuntime{
			game:    s.game,
			matchID: matchID,
		}

		s.matchRuntimes[matchID] = mRuntime
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
			delete(s.matchRuntimes, matchID)
		}
	}

	return ok
}

type matchRuntime struct {
	game Game

	matchID string
	guard   int
}

func (s *matchRuntime) reqContainerHandler(trn store.Transaction, reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
	match, err := trn.MatchRead(s.matchID)
	if err != nil {
		return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err, trn.Abort())}
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
			if err := trn.Commit(); err != nil {
				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return &protocol.ResponseContainer{Response: protocol.ResNotImplemented}
		}).
		On(protocol.ReqListRequests, func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			return &protocol.ResponseContainer{
				Response: protocol.ResOK,
				Body:     ctx.Requests(),
			}
		}).
		After(func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			if err := trn.Commit(); err != nil {
				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return nil
		})

	return ctx
}

func joinContext(trn store.Transaction, reqContainer *protocol.RequestContainer) *Context {
	return NewContext().
		On(protocol.ReqPlayerWantToJoin, func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			if err := trn.MatchPlayerJoins(reqContainer.Match.ID, reqContainer.Player.Name); err != nil {
				if err := trn.Abort(); err != nil {
					return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
				}

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
				if err := trn.Abort(); err != nil {
					return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
				}

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
				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err, trn.Abort())}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		})
}

func askToStartMatchContext(trn store.Transaction, reqContainer *protocol.RequestContainer) *Context {
	return NewContext().
		On(protocol.ReqPlayerAccepts, func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			msg := protocol.Message{Code: protocol.MPlayerAccepts}
			msg.EncodePayload(reqContainer.Player.Name)

			if err := trn.MatchAddMessageToCurrentDeal(reqContainer.Match.ID, msg); err != nil {
				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err, trn.Abort())}
			}

			match, err := trn.MatchRead(reqContainer.Match.ID)
			if err != nil {
				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err, trn.Abort())}
			}

			var OKPlayersName []string
			for _, m := range protocol.IndexDeals(match.Deals, 0).FindAllByMessageCode(protocol.MPlayerAccepts) {
				var mPayload string
				m.DecodePayload(&mPayload)

				OKPlayersName = append(OKPlayersName, mPayload)
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
						return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err, trn.Abort())}
					}

					return &protocol.ResponseContainer{Response: protocol.ResOK}
				}
			}

			if err := trn.MatchAddMessageToCurrentDeal(reqContainer.Match.ID, protocol.Message{Code: protocol.MDealCompleted}); err != nil {
				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err, trn.Abort())}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		}).
		On(protocol.ReqPlayerRefuses, func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			if err := trn.MatchAddMessageToCurrentDeal(reqContainer.Match.ID, protocol.Message{Code: protocol.MDealAborded}); err != nil {
				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err, trn.Abort())}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		})
}
