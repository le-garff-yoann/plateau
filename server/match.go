package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"plateau/protocol"
	"plateau/server/response"
	"plateau/server/response/body"
	"plateau/store"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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
		reqBody protocol.Match

		now = time.Now()

		err error
	)

	if err = json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, body.New().Ko(err))

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
			response.WriteJSON(w, http.StatusNotFound, body.New().Ko(fmt.Errorf(`Match "%s" not found`, v["id"])))
		} else {
			response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))
		}

		return
	}

	response.WriteJSON(w, http.StatusOK, match)
}

func (s *Server) getMatchConnectedPlayersNameHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	match, err := s.store.Matchs().Read(v["id"])
	if err != nil {
		if _, ok := err.(store.DontExistError); ok {
			response.WriteJSON(w, http.StatusNotFound, body.New().Ko(fmt.Errorf(`Match "%s" not found`, v["id"])))
		} else {
			response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))
		}

		return
	}

	var names []string
	for _, p := range match.ConnectedPlayers {
		names = append(names, p.Name)
	}

	response.WriteJSON(w, http.StatusOK, names)
}

func (s *Server) getMatchPlayersNameHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	match, err := s.store.Matchs().Read(v["id"])
	if err != nil {
		if _, ok := err.(store.DontExistError); ok {
			response.WriteJSON(w, http.StatusNotFound, body.New().Ko(fmt.Errorf(`Match "%s" not found`, v["id"])))
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

func (s *Server) getMatchTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	match, err := s.store.Matchs().Read(v["id"])
	if err != nil {
		if _, ok := err.(store.DontExistError); ok {
			response.WriteJSON(w, http.StatusNotFound, body.New().Ko(fmt.Errorf(`Match "%s" not found`, v["id"])))
		} else {
			response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))
		}

		return
	}

	response.WriteJSON(w, http.StatusOK, match.Transactions)
}

func (s *Server) connectMatchHandler(w http.ResponseWriter, r *http.Request) {
	session, err := s.store.Sessions().Get(r, ServerName)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	var (
		username = session.Values["username"].(string)

		v = mux.Vars(r)
	)

	done := make(chan int)

	s.doneWg.Add(1)
	defer s.doneWg.Done()

	srvDoneCh, srvDoneUUID := s.doneBroadcaster.Subscribe()
	defer s.doneBroadcaster.Unsubscribe(srvDoneUUID)

	if err := s.store.Matchs().ConnectPlayer(v["id"], username); err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(store.PlayerConnectionError); ok {
			statusCode = http.StatusConflict
		} else {
			logrus.Error(err)
		}

		w.WriteHeader(statusCode)

		return
	}
	defer s.store.Matchs().DisconnectPlayer(v["id"], username)

	mRuntime, err := s.guardRuntime(v["id"])
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	defer s.unguardRuntime(v["id"])

	c, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	defer c.Close()

	go func() {
		for {
			select {
			case <-done:
				return
			case <-srvDoneCh:
				s.unguardRuntime(v["id"])

				s.store.Matchs().DisconnectPlayer(v["id"], username)

				s.doneBroadcaster.Unsubscribe(srvDoneUUID)
				s.doneWg.Done()

				return
			}
		}
	}()
	defer func() {
		done <- 0
	}()

	notificationCh, notificationUUID := mRuntime.transactionsChangesBroadcaster.Subscribe()
	defer mRuntime.transactionsChangesBroadcaster.Unsubscribe(notificationUUID)

	var writeJSONMux sync.Mutex
	writeJSON := func(v interface{}) error {
		writeJSONMux.Lock()
		defer writeJSONMux.Unlock()

		return c.WriteJSON(v)
	}

	logCtx := logrus.
		WithField("match", v["id"]).
		WithField("player", username)

	go func() {
		for {
			v, ok := <-notificationCh
			if !ok {
				return
			}

			notificationContainer := v.(protocol.NotificationContainer)

			logCtx.
				WithField("type", "notification").
				Debug(notificationContainer)

			writeJSON(notificationContainer)
		}
	}()

	for {
		var requestContainer protocol.RequestContainer

		if err := c.ReadJSON(&requestContainer); err != nil {
			if writeJSON(protocol.ResponseContainer{Response: protocol.ResBadRequest, Body: body.New().Ko(err)}) != nil {
				return
			}

			continue
		}

		var err error
		requestContainer.Player, err = s.store.Players().Read(username)
		if err != nil {
			if writeJSON(protocol.ResponseContainer{Response: protocol.ResInternalError, Body: body.New().Ko(err)}) != nil {
				return
			}

			continue
		}

		logCtx.
			WithField("type", "request").
			Debug(requestContainer)

		responseContainer := mRuntime.requestContainerHandler(&requestContainer)

		logCtx.
			WithField("type", "response").
			Debug(responseContainer)

		writeJSON(responseContainer)
	}
}
