package server

import (
	"net/http"
	"plateau/model"
	"plateau/server/response"
	"plateau/server/response/body"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func (s *Server) getPlayersName(w http.ResponseWriter, r *http.Request) {
	var names []string

	if errs := s.db.Order("name").Select("name").Find(&[]model.Player{}).Pluck("name", &names).GetErrors(); len(errs) > 0 {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(errs...))

		return
	}

	response.WriteJSON(w, http.StatusOK, names)
}

func (s *Server) readPlayer(w http.ResponseWriter, r *http.Request) {
	var (
		v = mux.Vars(r)

		player model.Player
	)

	if errs := s.db.Where("name = ?", v["name"]).First(&player).GetErrors(); len(errs) > 0 {
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

	response.WriteJSON(w, http.StatusOK, player)
}

func (s *Server) getPlayerGamesID(w http.ResponseWriter, r *http.Request) {
	var (
		v = mux.Vars(r)

		player model.Player
	)

	if errs := s.db.Preload("Games", func(db *gorm.DB) *gorm.DB {
		return db.Order("id")
	}).Where("name = ?", v["name"]).First(&player).GetErrors(); len(errs) > 0 {
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

	var gamesID []uint
	for _, g := range player.Games {
		gamesID = append(gamesID, g.ID)
	}

	response.WriteJSON(w, http.StatusOK, gamesID)
}
