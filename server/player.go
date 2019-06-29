package server

import (
	"fmt"
	"net/http"
	"plateau/server/response"
	"plateau/server/response/body"
	"plateau/store"

	"github.com/gorilla/mux"
)

func (s *Server) getPlayersNameHandler(w http.ResponseWriter, r *http.Request) {
	trn := s.store.BeginTransaction()
	defer trn.Abort()

	names, err := trn.PlayerList()
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	response.WriteJSON(w, http.StatusOK, names)
}

func (s *Server) readPlayerHandler(w http.ResponseWriter, r *http.Request) {
	playerName := mux.Vars(r)["name"]

	trn := s.store.BeginTransaction()
	defer trn.Abort()

	player, err := trn.PlayerRead(playerName)
	if err != nil {
		if _, ok := err.(store.DontExistError); ok {
			response.WriteJSON(w, http.StatusNotFound, body.New().Ko(fmt.Errorf(`Player "%s" not found`, playerName)))
		} else {
			response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))
		}

		return
	}

	response.WriteJSON(w, http.StatusOK, player)
}
