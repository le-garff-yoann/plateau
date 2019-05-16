package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"plateau/game"
	"plateau/server/response"
	"plateau/server/response/body"
	"plateau/store"
	"time"

	"github.com/gorilla/mux"
)

func (s *Server) getMatchIDsHandler(w http.ResponseWriter, r *http.Request) {
	IDs, err := s.store.Matchs().List()
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	response.WriteJSON(w, http.StatusOK, IDs)
}

func (s *Server) createMatchHandler(w http.ResponseWriter, r *http.Request) {
	var (
		reqBody store.Match
		g       = game.New()

		now = time.Now()

		err error
	)

	if err = g.Init(); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	if err = json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, body.New().Ko(err))

		return
	}

	if err = g.IsMatchValid(&reqBody); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, body.New().Ko(err))

		return
	}

	match := store.Match{
		CreatedAt:               &now,
		NumberOfPlayersRequired: reqBody.NumberOfPlayersRequired,
	}

	if match.ID, err = s.store.Matchs().Create(match); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	response.WriteJSON(w, http.StatusCreated, match)
}

func (s *Server) readMatchHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	match, err := s.store.Matchs().Read(v["id"])
	if err != nil {
		if _, ok := err.(store.DontExistError); ok {
			response.WriteJSON(w, http.StatusNotFound, body.New().Ko(fmt.Errorf("Match %s not found", v["id"])))
		} else {
			response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))
		}

		return
	}

	response.WriteJSON(w, http.StatusOK, match)
}

func (s *Server) getMatchPlayersNameHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	match, err := s.store.Matchs().Read(v["id"])
	if err != nil {
		if _, ok := err.(store.DontExistError); ok {
			response.WriteJSON(w, http.StatusNotFound, body.New().Ko(fmt.Errorf("Match %s not found", v["id"])))
		} else {
			response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))
		}

		return
	}

	var names []string
	for _, p := range match.Players {
		names = append(names, p.Name)
	}

	response.WriteJSON(w, http.StatusOK, names)
}

func (s *Server) getMatchEventContainersHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	match, err := s.store.Matchs().Read(v["id"])
	if err != nil {
		if _, ok := err.(store.DontExistError); ok {
			response.WriteJSON(w, http.StatusNotFound, body.New().Ko(fmt.Errorf("Match %s not found", v["id"])))
		} else {
			response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))
		}

		return
	}

	response.WriteJSON(w, http.StatusOK, match.EventContainers)
}
