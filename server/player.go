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
	names, err := s.store.Players().List()
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	response.WriteJSON(w, http.StatusOK, names)
}

func (s *Server) readPlayerHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	player, err := s.store.Players().Read(v["name"])
	if err != nil {
		if _, ok := err.(store.DontExistError); ok {
			response.WriteJSON(w, http.StatusNotFound, body.New().Ko(fmt.Errorf("Player %s not found", v["name"])))
		} else {
			response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))
		}

		return
	}

	response.WriteJSON(w, http.StatusOK, player)
}
