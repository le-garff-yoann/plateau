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
	trn, err := s.store.BeginTransaction()
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	names, err := trn.PlayerList()
	if err := trn.Abort(); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	if names == nil {
		names = []string{}
	}

	response.WriteJSON(w, http.StatusOK, names)
}

func (s *Server) readPlayerHandler(w http.ResponseWriter, r *http.Request) {
	playerName := mux.Vars(r)["name"]

	trn, err := s.store.BeginTransaction()
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	player, err := trn.PlayerRead(playerName)
	if err := trn.Abort(); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

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
