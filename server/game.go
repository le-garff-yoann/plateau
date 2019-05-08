package server

import (
	"encoding/json"
	"log"
	"net/http"
	"plateau/game"
	"plateau/model"
	"plateau/server/response"
	"plateau/server/response/body"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func (s *Server) getGameIDs(w http.ResponseWriter, r *http.Request) {
	var IDs []uint

	if errs := s.db.Order("id").Select("id").Find(&[]model.Game{}).Pluck("id", &IDs).GetErrors(); len(errs) > 0 {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(errs...))

		return
	}

	response.WriteJSON(w, http.StatusOK, IDs)
}

func (s *Server) createGame(w http.ResponseWriter, r *http.Request) {
	var (
		reqBody model.Game

		now = time.Now()
	)

	game, err := game.New()
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, body.New().Ko(err))

		return
	}

	if err := game.IsValid(&reqBody); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, body.New().Ko(err))

		return
	}

	gameModel := model.Game{
		CreatedAt:               &now,
		NumberOfPlayersRequired: reqBody.NumberOfPlayersRequired,
	}

	if errs := s.db.Create(&gameModel).GetErrors(); len(errs) > 0 {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(errs...))

		return
	}

	response.WriteJSON(w, http.StatusCreated, gameModel)
}

func (s *Server) readGame(w http.ResponseWriter, r *http.Request) {
	var (
		v = mux.Vars(r)

		gameModel model.Game
	)

	if errs := s.db.First(&gameModel, v["id"]).GetErrors(); len(errs) > 0 {
		httpCode := http.StatusInternalServerError

		for _, err := range errs {
			if gorm.IsRecordNotFoundError(err) {
				httpCode = http.StatusNotFound

				break
			}
		}

		response.WriteJSON(w, httpCode, body.New().Ko(errs...))

		return
	}

	response.WriteJSON(w, http.StatusOK, gameModel)
}

func (s *Server) connectGame(w http.ResponseWriter, r *http.Request) {
	session, err := s.sessionStore.Get(r, ServerName)
	if err != nil {
		return
	}

	v := mux.Vars(r)

	c, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer func() {
		// TODO: Sent `model.EPlayerDisconnects` to everyone else (we would need an "events store" built on top of the store).

		c.Close()
	}()

	game, err := game.New()
	if err != nil {
		return
	}

	// TODO: Sent `model.EPlayerConnects` to everyone else (we would need an "events store" built on top of the store).
	// eventContainer := model.EventContainer{Event: model.EPlayerConnects, Subjects: []string{playerName}}

	// if c.WriteJSON(eventContainer) != nil {
	// 	return
	// }

	// log.Printf(eventContainer.String())

	// REFACTOR: Most of this function's code SHOULD BE in another package ("events store"). See TODOs.

	for {
		var eventContainer model.EventContainer

		if c.ReadJSON(&eventContainer) != nil {
			if c.WriteJSON(model.EventContainer{Event: model.EIllegal}) != nil {
				return
			}

			continue
		}

		var (
			gameModel model.Game

			playerModel model.Player
			isPlayerIn  = false
		)

		if errs := s.db.Preload("Players").First(&gameModel, v["id"]).GetErrors(); len(errs) > 0 {
			if c.WriteJSON(model.EventContainer{Event: model.EInternalError}) != nil {
				return
			}

			continue
		}

		if errs := s.db.Where("name = ?", session.Values["username"].(string)).First(&playerModel).GetErrors(); len(errs) > 0 {
			if c.WriteJSON(model.EventContainer{Event: model.EInternalError}) != nil {
				return
			}

			continue
		}

		eventContainer.Emitter = &playerModel

		for _, p := range gameModel.Players {
			if p.Name == eventContainer.Emitter.Name {
				isPlayerIn = true

				break
			}
		}

		eventContainer.Game = &gameModel

		log.Printf(eventContainer.String())

		if !eventContainer.IsLegal() {
			if c.WriteJSON(model.EventContainer{Event: model.EIllegal}) != nil {
				return
			}

			continue
		}

		if game.OnEvent(&eventContainer) != nil {
			if c.WriteJSON(model.EventContainer{Event: model.EIllegal}) != nil {
				return
			}

			continue
		}

		switch eventContainer.Event {
		case model.EPlayerWantToJoin:
			if isPlayerIn {
				if c.WriteJSON(model.EventContainer{Event: model.EIllegal}) != nil {
					return
				}

				continue
			}

			if err := s.db.Model(&playerModel).Association("Games").Append(gameModel).Error; err != nil {
				if c.WriteJSON(model.EventContainer{Event: model.EInternalError}) != nil {
					return
				}

				continue
			}

			eventContainer.Subjects = append(eventContainer.Subjects, eventContainer.Emitter)

			pairedEventContainer := &eventContainer
			pairedEventContainer.Event = model.EPlayerJoins

			game.OnEvent(pairedEventContainer)
		case model.EPlayerWantToLeave:
			if !isPlayerIn || gameModel.Running {
				if c.WriteJSON(model.EventContainer{Event: model.EIllegal}) != nil {
					return
				}

				continue
			}

			if s.db.Model(&playerModel).Association("Games").Delete(gameModel).Error != nil {
				if c.WriteJSON(model.EventContainer{Event: model.EInternalError}) != nil {
					return
				}

				continue
			}

			eventContainer.Subjects = append(eventContainer.Subjects, eventContainer.Emitter)

			pairedEventContainer := &eventContainer
			pairedEventContainer.Event = model.EPlayerLeaves

			game.OnEvent(pairedEventContainer)
		}

		if c.WriteJSON(model.EventContainer{Event: model.EProcessed}) != nil {
			return
		}
	}
}

func (s *Server) getGamePlayersName(w http.ResponseWriter, r *http.Request) {
	var (
		v = mux.Vars(r)

		game model.Game
	)

	if errs := s.db.Preload("Players", func(db *gorm.DB) *gorm.DB {
		return db.Order("name")
	}).First(&game, v["id"]).GetErrors(); len(errs) > 0 {
		httpCode := http.StatusInternalServerError

		for _, err := range errs {
			if gorm.IsRecordNotFoundError(err) {
				httpCode = http.StatusNotFound

				break
			}
		}

		response.WriteJSON(w, httpCode, body.New().Ko(errs...))

		return
	}

	var playersName []string
	for _, p := range game.Players {
		playersName = append(playersName, p.Name)
	}

	response.WriteJSON(w, http.StatusOK, playersName)
}

func (s *Server) getGameEventContainers(w http.ResponseWriter, r *http.Request) {
	var (
		v = mux.Vars(r)

		game model.Game
	)

	if errs := s.db.Preload("EventContainers", func(db *gorm.DB) *gorm.DB {
		return db.Order("id")
	}).First(&game, v["id"]).GetErrors(); len(errs) > 0 {
		httpCode := http.StatusInternalServerError

		for _, err := range errs {
			if gorm.IsRecordNotFoundError(err) {
				httpCode = http.StatusNotFound

				break
			}
		}

		response.WriteJSON(w, httpCode, body.New().Ko(errs...))

		return
	}

	response.WriteJSON(w, http.StatusOK, game.EventContainers)
}
