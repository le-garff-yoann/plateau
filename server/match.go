package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"plateau/protocol"
	"plateau/server/response"
	"plateau/server/response/body"
	"plateau/store"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (s *Server) getMatchIDsHandler(w http.ResponseWriter, r *http.Request) {
	trn := s.store.BeginTransaction()

	IDs, err := trn.MatchList()
	trn.Abort()

	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	if IDs == nil {
		IDs = []string{}
	}

	response.WriteJSON(w, http.StatusOK, IDs)
}

func (s *Server) createMatchHandler(w http.ResponseWriter, r *http.Request) {
	var (
		reqBody protocol.Match

		now = time.Now()

		err error
	)

	if err = json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, body.New().Ko(err))

		return
	}

	if !(reqBody.NumberOfPlayersRequired >= s.game.MinPlayers() && reqBody.NumberOfPlayersRequired <= s.game.MaxPlayers()) {
		response.WriteJSON(w, http.StatusBadRequest, body.New().Ko(
			fmt.Errorf("The number of players must be between %d and %d", s.game.MinPlayers(), s.game.MaxPlayers()),
		))

		return
	}
	if err = s.game.IsMatchValid(&reqBody); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, body.New().Ko(err))

		return
	}

	match := protocol.Match{
		CreatedAt:               now,
		NumberOfPlayersRequired: reqBody.NumberOfPlayersRequired,
	}

	trn := s.store.BeginTransaction()

	if match.ID, err = trn.MatchCreate(match); err != nil {
		trn.Abort()

		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	trn.Commit()

	response.WriteJSON(w, http.StatusCreated, match)
}

func (s *Server) readMatchHandler(w http.ResponseWriter, r *http.Request) {
	matchID := mux.Vars(r)["id"]

	trn := s.store.BeginTransaction()

	match, err := trn.MatchRead(matchID)
	trn.Abort()

	if err != nil {
		if _, ok := err.(store.DontExistError); ok {
			response.WriteJSON(w, http.StatusNotFound, body.New().Ko(fmt.Errorf(`Match "%s" not found`, matchID)))
		} else {
			response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))
		}

		return
	}

	response.WriteJSON(w, http.StatusOK, match)
}

func (s *Server) getMatchPlayersNameHandler(w http.ResponseWriter, r *http.Request) {
	matchID := mux.Vars(r)["id"]

	trn := s.store.BeginTransaction()

	match, err := trn.MatchRead(matchID)
	trn.Abort()

	if err != nil {
		if _, ok := err.(store.DontExistError); ok {
			response.WriteJSON(w, http.StatusNotFound, body.New().Ko(fmt.Errorf(`Match "%s" not found`, matchID)))
		} else {
			response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))
		}

		return
	}

	names := []string{}
	for _, p := range match.Players {
		names = append(names, p.Name)
	}

	response.WriteJSON(w, http.StatusOK, names)
}

func (s *Server) streamMatchNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(errors.New("Cannot flush")))

		return
	}

	var (
		matchID = mux.Vars(r)["id"]

		trn = s.store.BeginTransaction()
	)

	match, err := trn.MatchRead(matchID)
	trn.Abort()

	if err != nil {
		if _, ok := err.(store.DontExistError); ok {
			response.WriteJSON(w, http.StatusNotFound, body.New().Ko(fmt.Errorf(`Match "%s" not found`, matchID)))
		} else {
			response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))
		}

		return
	}

	if match.EndedAt != nil {
		response.WriteJSON(w, http.StatusGone, body.New().Ko(fmt.Errorf(`Match "%s" is ended`, match.ID)))

		return
	}

	ch := make(chan interface{})
	s.store.RegisterNotificationsChannel(ch)
	defer s.store.UnregisterNotificationsChannel(ch)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	w.WriteHeader(http.StatusOK)

	flusher.Flush()

	for {
		v, ok := <-ch
		if !ok {
			return
		}

		matchNotification, ok := v.(store.MatchNotification)
		if ok && matchNotification.ID == match.ID {
			logCtx := logrus.
				WithField("match", match.ID)

			logCtx.Debug(matchNotification)

			jsonb, err := json.Marshal(matchNotification)
			if err != nil {
				logCtx.Panic(err)
			}

			_, err = fmt.Fprintf(w, "data: %s\n\n", string(jsonb))
			if err != nil {
				logCtx.Warning(err)

				return
			}

			flusher.Flush()
		}
	}
}

func (s *Server) getMatchDealsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		trn = s.store.BeginTransaction()

		matchID = mux.Vars(r)["id"]
	)

	match, err := trn.MatchRead(matchID)
	trn.Abort()

	if err != nil {
		if _, ok := err.(store.DontExistError); ok {
			response.WriteJSON(w, http.StatusNotFound, body.New().Ko(fmt.Errorf(`Match "%s" not found`, matchID)))
		} else {
			response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))
		}

		return
	}

	var playerName []string
	if !match.IsEnded() {
		session, err := s.store.Sessions().Get(r, ServerName)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

			return
		}

		playerName = append(playerName, session.Values["username"].(string))
	}

	deals := []protocol.Deal{}
	for _, d := range match.Deals {
		deals = append(deals, *d.WithMessagesConcealed(playerName...))
	}

	response.WriteJSON(w, http.StatusOK, deals)
}

func (s *Server) patchMatchHandler(w http.ResponseWriter, r *http.Request) {
	session, err := s.store.Sessions().Get(r, ServerName)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	var (
		username = session.Values["username"].(string)
		matchID  = mux.Vars(r)["id"]
	)

	s.doneWg.Add(1)
	defer s.doneWg.Done()

	mRuntime, err := s.guardRuntime(matchID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}
	defer s.unguardRuntime(matchID)

	logCtx := logrus.
		WithField("match", matchID).
		WithField("player", username)

	var reqContainer protocol.RequestContainer
	if err = json.NewDecoder(r.Body).Decode(&reqContainer); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, body.New().Ko(err))

		return
	}

	trn := s.store.BeginTransaction()

	reqContainer.Player, err = trn.PlayerRead(username)
	if err != nil {
		trn.Abort()

		response.WriteJSON(w, http.StatusBadRequest, body.New().Ko(err))

		return
	}

	logCtx.
		WithField("type", "request").
		Debug(reqContainer)

	resContainer := mRuntime.reqContainerHandler(trn, &reqContainer)
	if !trn.Closed() {
		logCtx.Panic("Transaction not closed")
	}

	logCtx.
		WithField("type", "response").
		Debug(resContainer)

	response.WriteJSON(w, http.StatusOK, resContainer)
}
