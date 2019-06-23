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
		iterator, err := s.store.Matchs().CreateTransactionsChangeIterator(matchID)
		if err != nil {
			return nil, err
		}

		s.matchmatchRuntimesMux.Lock()
		defer s.matchmatchRuntimesMux.Unlock()
		s.matchRuntimes[matchID] = &MatchRuntime{
			game:                           s.game,
			matchID:                        matchID,
			matchs:                         s.store.Matchs(),
			Matchs:                         s.store.Matchs(),
			Players:                        s.store.Players(),
			transactionsChanges:            []store.TransactionChange{},
			transactionsChangesBroadcaster: broadcaster.New(),
			done:                           make(chan int),
		}

		go s.matchRuntimes[matchID].transactionsChangesBroadcaster.Run()

		go func() {
			var trxChange store.TransactionChange

			for iterator.Next(&trxChange) {
				func(r *MatchRuntime) {
					r.transactionsChangesMux.Lock()
					defer r.transactionsChangesMux.Unlock()

					r.transactionsChanges = append(r.transactionsChanges, trxChange)

					r.transactionsChangesBroadcaster.Submit(protocol.NotificationContainer{
						Notification: protocol.NTransactionChange,
						Body:         trxChange,
					})
				}(s.matchRuntimes[matchID])
			}
		}()

		go func(r *MatchRuntime) {
			<-r.done

			r.transactionsChangesBroadcaster.Done()

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

	matchs  store.MatchStore
	Matchs  store.MatchGameStore
	Players store.PlayerGameStore

	transactionsChangesMux         sync.RWMutex
	transactionsChanges            []store.TransactionChange
	transactionsChangesBroadcaster *broadcaster.Broadcaster

	done chan int
}

// TransactionsChanges ...
func (s *MatchRuntime) TransactionsChanges() []store.TransactionChange {
	s.transactionsChangesMux.RLock()
	defer s.transactionsChangesMux.RUnlock()

	return s.transactionsChanges
}

func (s *MatchRuntime) requestContainerHandler(requestContainer *protocol.RequestContainer) *protocol.ResponseContainer {
	match, err := s.Matchs.Read(s.matchID)
	if err != nil {
		return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
	}

	requestContainer.Match = match

	ctx := baseContext(requestContainer)

	if requestContainer.Match.EndedAt == nil {
		for _, trx := range requestContainer.Match.Transactions {
			if trx.FindByMessageCode(protocol.MPlayerWantToStartTheGame) != nil && trx.FindByMessageCode(protocol.MTransactionCompleted) != nil {
				return baseContext(requestContainer).Complete(
					s.game.Context(s, requestContainer),
				).handle(s, requestContainer)
			}
		}

		currentTrx := protocol.IndexTransactions(requestContainer.Match.Transactions, 0)

		if currentTrx == nil || !currentTrx.IsActive() {
			if funk.Contains(requestContainer.Match.Players, *requestContainer.Player) {
				ctx.Complete(leaveContext(requestContainer)).Complete(wantToStartMatchContext(requestContainer))
			} else {
				ctx.Complete(joinContext(requestContainer))
			}
		} else {
			if currentTrx.Holder.Name == requestContainer.Player.Name {
				ctx.Complete(askToStartMatchContext(requestContainer))
			}
		}
	}

	return ctx.handle(s, requestContainer)
}

func baseContext(requestContainer *protocol.RequestContainer) *Context {
	ctx := NewContext()
	ctx.
		On(protocol.ReqGetCurrentTransaction, func(matchRuntime *MatchRuntime, requestContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			return &protocol.ResponseContainer{
				Response: protocol.ResOK,
				Body:     protocol.IndexTransactions(requestContainer.Match.Transactions, 0),
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

func joinContext(requestContainer *protocol.RequestContainer) *Context {
	return NewContext().
		On(protocol.ReqPlayerWantToJoin, func(matchRuntime *MatchRuntime, requestContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			if err := matchRuntime.matchs.PlayerJoins(matchRuntime.matchID, requestContainer.Player.Name); err != nil {
				if _, ok := err.(store.PlayerParticipationError); ok {
					return &protocol.ResponseContainer{Response: protocol.ResForbidden, Body: body.New().Ko(err)}
				}

				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		})
}

func leaveContext(requestContainer *protocol.RequestContainer) *Context {
	return NewContext().
		On(protocol.ReqPlayerWantToLeave, func(matchRuntime *MatchRuntime, requestContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			if err := matchRuntime.matchs.PlayerLeaves(matchRuntime.matchID, requestContainer.Player.Name); err != nil {
				if _, ok := err.(store.PlayerParticipationError); ok {
					return &protocol.ResponseContainer{Response: protocol.ResForbidden, Body: body.New().Ko(err)}
				}

				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		})
}

func wantToStartMatchContext(requestContainer *protocol.RequestContainer) *Context {
	return NewContext().
		On(protocol.ReqPlayerWantToStartTheGame, func(matchRuntime *MatchRuntime, requestContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			if !requestContainer.Match.IsFull() {
				return &protocol.ResponseContainer{
					Response: protocol.ResForbidden,
					Body:     body.New().Ko(errors.New("There are not enough players")),
				}
			}

			if err := matchRuntime.Matchs.CreateTransaction(matchRuntime.matchID, protocol.Transaction{Holder: *requestContainer.Player, Messages: []protocol.Message{protocol.Message{MessageCode: protocol.MPlayerWantToStartTheGame}}}); err != nil {
				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		})
}

func askToStartMatchContext(requestContainer *protocol.RequestContainer) *Context {
	return NewContext().
		On(protocol.ReqPlayerAccepts, func(matchRuntime *MatchRuntime, requestContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			// REFACTOR: Should defer the message deletion in the case of a subsequent failure in this routine.
			if err := matchRuntime.Matchs.AddMessageToCurrentTransaction(matchRuntime.matchID, protocol.Message{MessageCode: protocol.MPlayerAccepts, Payload: requestContainer.Player.Name}); err != nil {
				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			match, err := matchRuntime.Matchs.Read(matchRuntime.matchID)
			if err != nil {
				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			var OKPlayersName []string
			for _, m := range protocol.IndexTransactions(match.Transactions, 0).FindAllByMessageCode(protocol.MPlayerAccepts) {
				OKPlayersName = append(OKPlayersName, m.Payload.(string))
			}

			for _, player := range match.Players {
				if !funk.Contains(OKPlayersName, player.Name) {
					if err := matchRuntime.Matchs.UpdateCurrentTransactionHolder(matchRuntime.matchID, player.Name); err != nil {
						return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
					}

					return &protocol.ResponseContainer{Response: protocol.ResOK}
				}
			}

			if err := matchRuntime.Matchs.AddMessageToCurrentTransaction(matchRuntime.matchID, protocol.Message{MessageCode: protocol.MTransactionCompleted}); err != nil {
				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		}).
		On(protocol.ReqPlayerRefuses, func(matchRuntime *MatchRuntime, requestContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			if err := matchRuntime.Matchs.AddMessageToCurrentTransaction(matchRuntime.matchID, protocol.Message{MessageCode: protocol.MTransactionAborded}); err != nil {
				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return &protocol.ResponseContainer{Response: protocol.ResOK}
		})
}
