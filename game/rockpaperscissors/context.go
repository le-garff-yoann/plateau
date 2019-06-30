package rockpaperscissors

import (
	"plateau/protocol"
	"plateau/server"
	"plateau/server/response/body"
	"plateau/store"
	"time"
)

// Context implements `server.Game` interface.
func (s *Game) Context(trn store.Transaction, reqContainer *protocol.RequestContainer) *server.Context {
	return server.NewContext().
		Before(func(ctx *server.Context, reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			currentDeal := protocol.IndexDeals(reqContainer.Match.Deals, 0)

			if (currentDeal.FindByMessageCode(protocol.MPlayerWantToStartTheGame) != nil && currentDeal.FindByMessageCode(protocol.MDealCompleted) != nil) ||
				currentDeal.FindByMessageCode(protocol.MDealAborded) != nil {

				err := trn.MatchCreateDeal(reqContainer.Match.ID, protocol.Deal{Holder: *reqContainer.Match.RandomPlayer()})
				if err != nil {
					trn.Abort()

					return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
				}

				reqContainer.Match, err = trn.MatchRead(reqContainer.Match.ID)
				if err != nil {
					trn.Abort()

					return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
				}
			}

			if reqContainer.Player.Name == protocol.IndexDeals(reqContainer.Match.Deals, 0).Holder.Name {
				ctx.
					On(ReqRock, requestFunc(MRock, trn)).
					On(ReqPaper, requestFunc(MPaper, trn)).
					On(ReqScissors, requestFunc(MScissors, trn))
			}

			return nil
		}).
		After(func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			match, err := trn.MatchRead(reqContainer.Match.ID)
			if err != nil {
				trn.Abort()

				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			currentDeal := protocol.IndexDeals(match.Deals, 0)

			if len(currentDeal.Messages) >= len(match.Players) {
				var (
					playerOneMessage = currentDeal.Find(func(msg protocol.Message) bool { return msg.Payload.(string) == match.Players[0].Name })
					playerTwoMessage = currentDeal.Find(func(msg protocol.Message) bool { return msg.Payload.(string) == match.Players[1].Name })
				)

				if playerOneMessage.MessageCode == playerTwoMessage.MessageCode {
					if err := trn.MatchAddMessageToCurrentDeal(match.ID, protocol.Message{MessageCode: protocol.MDealAborded}); err != nil {
						trn.Abort()

						return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
					}
				} else {
					var winner, loser protocol.Player

					if (playerOneMessage.MessageCode == MPaper && playerTwoMessage.MessageCode == MRock) ||
						(playerOneMessage.MessageCode == MScissors && playerTwoMessage.MessageCode == MPaper) ||
						(playerOneMessage.MessageCode == MRock && playerTwoMessage.MessageCode == MScissors) {
						winner, loser = match.Players[0], match.Players[1]
					} else {
						winner, loser = match.Players[1], match.Players[0]
					}

					if err := trn.MatchAddMessageToCurrentDeal(match.ID, protocol.Message{MessageCode: protocol.MDealCompleted}); err != nil {
						trn.Abort()

						return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
					}

					if err := trn.PlayerIncreaseWins(winner.Name, 1); err != nil {
						trn.Abort()

						return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
					}
					if err := trn.PlayerIncreaseLoses(loser.Name, 1); err != nil {
						trn.Abort()

						return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
					}

					if err := trn.MatchEndedAt(match.ID, time.Now()); err != nil {
						trn.Abort()

						return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
					}
				}
			}

			trn.Commit()

			return nil
		})
}

func requestFunc(msg protocol.MessageCode, trn store.Transaction) func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
	return func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
		if err := trn.MatchAddMessageToCurrentDeal(reqContainer.Match.ID, protocol.Message{
			MessageCode: msg,
			Payload:     reqContainer.Player.Name,
		}); err != nil {
			trn.Abort()

			return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
		}

		if err := trn.MatchUpdateCurrentDealHolder(reqContainer.Match.ID, reqContainer.Match.NextPLayer(*reqContainer.Player).Name); err != nil {
			trn.Abort()

			return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
		}

		return &protocol.ResponseContainer{Response: protocol.ResOK}
	}
}
