package rockpaperscissors

import (
	"plateau/protocol"
	"plateau/server"
	"plateau/server/response/body"
	"plateau/store"
	"time"
)

// Context implements the `server.Game` interface.
func (s *Game) Context(trn store.Transaction, reqContainer *protocol.RequestContainer) *server.Context {
	return server.NewContext().
		Before(func(ctx *server.Context, reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
			currentDeal := protocol.IndexDeals(reqContainer.Match.Deals, 0)

			if (currentDeal.FindByMessageCode(protocol.MPlayerWantToStartTheGame) != nil && currentDeal.FindByMessageCode(protocol.MDealCompleted) != nil) ||
				currentDeal.FindByMessageCode(protocol.MDealAborded) != nil {

				err := trn.MatchCreateDeal(reqContainer.Match.ID, protocol.Deal{Holder: *reqContainer.Match.RandomPlayer()})
				if err != nil {
					return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err, trn.Abort())}
				}

				reqContainer.Match, err = trn.MatchRead(reqContainer.Match.ID)
				if err != nil {
					return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err, trn.Abort())}
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
				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err, trn.Abort())}
			}

			currentDeal := protocol.IndexDeals(match.Deals, 0)

			if len(currentDeal.Messages) >= len(match.Players) {
				var (
					playerAMessage = currentDeal.Find(func(msg protocol.Message) bool {
						var msgPayload string
						msg.DecodePayload(&msgPayload)

						return msgPayload == match.Players[0].Name
					})

					playerBMessage = currentDeal.Find(func(msg protocol.Message) bool {
						var msgPayload string
						msg.DecodePayload(&msgPayload)

						return msgPayload == match.Players[1].Name
					})
				)

				if playerAMessage.Code == playerBMessage.Code {
					if err := trn.MatchAddMessageToCurrentDeal(match.ID, protocol.Message{Code: protocol.MDealAborded}); err != nil {
						return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err, trn.Abort())}
					}
				} else {
					var winner, loser protocol.Player

					if (playerAMessage.Code == MPaper && playerBMessage.Code == MRock) ||
						(playerAMessage.Code == MScissors && playerBMessage.Code == MPaper) ||
						(playerAMessage.Code == MRock && playerBMessage.Code == MScissors) {
						winner, loser = match.Players[0], match.Players[1]
					} else {
						winner, loser = match.Players[1], match.Players[0]
					}

					if err := trn.MatchAddMessageToCurrentDeal(match.ID, protocol.Message{Code: protocol.MDealCompleted}); err != nil {
						return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err, trn.Abort())}
					}

					if err := trn.PlayerIncreaseWins(winner.Name, 1); err != nil {
						return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err, trn.Abort())}
					}
					if err := trn.PlayerIncreaseLoses(loser.Name, 1); err != nil {
						return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err, trn.Abort())}
					}

					if err := trn.MatchEndedAt(match.ID, time.Now()); err != nil {
						return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err, trn.Abort())}
					}
				}
			}

			if err := trn.Commit(); err != nil {
				return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}
			}

			return nil
		})
}

func requestFunc(msg protocol.MessageCode, trn store.Transaction) func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
	return func(reqContainer *protocol.RequestContainer) *protocol.ResponseContainer {
		msg := protocol.Message{
			Code:             msg,
			AllowedNamesCode: []string{reqContainer.Player.Name},
		}
		msg.EncodePayload(reqContainer.Player.Name)

		if err := trn.MatchAddMessageToCurrentDeal(reqContainer.Match.ID, msg); err != nil {
			return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err, trn.Abort())}
		}

		if err := trn.MatchUpdateCurrentDealHolder(reqContainer.Match.ID, reqContainer.Match.NextPlayer(*reqContainer.Player).Name); err != nil {
			return &protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err, trn.Abort())}
		}

		return &protocol.ResponseContainer{Response: protocol.ResOK}
	}
}
