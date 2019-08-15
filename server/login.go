package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"plateau/protocol"
	"plateau/server/response"
	"plateau/server/response/body"
	"plateau/store"

	"golang.org/x/crypto/bcrypt"
)

type loginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *loginCredentials) IsValid() bool {
	return len(s.Username) > 0 && len(s.Password) > 0
}

func (s *loginCredentials) PasswordHash() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
}

func (s *loginCredentials) VerifyHash(h []byte) bool {
	return bcrypt.CompareHashAndPassword(h, []byte(s.Password)) == nil
}

func (s *Server) loginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, ServerName)
		if err != nil {
			response.WriteJSON(w, http.StatusForbidden, body.New().Ko(err))

			return
		}

		if session.Values["authenticated"] == nil || !session.Values["authenticated"].(bool) {
			w.WriteHeader(http.StatusForbidden)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var cred loginCredentials

	if json.NewDecoder(r.Body).Decode(&cred) != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if !cred.IsValid() {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	hPassword, err := cred.PasswordHash()
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	player := protocol.Player{Name: cred.Username, Password: string(hPassword)}

	trn, err := s.store.BeginTransaction()
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	if err := trn.PlayerCreate(player); err != nil {
		trn.Abort()

		if _, ok := err.(store.DuplicateError); ok {
			response.WriteJSON(w, http.StatusConflict, body.New().Ko(fmt.Errorf(`Player "%s" already exists`, cred.Username)))
		} else {
			response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))
		}

		return
	}

	trn.Commit()

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var cred loginCredentials

	if json.NewDecoder(r.Body).Decode(&cred) != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if !cred.IsValid() {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	trn, err := s.store.BeginTransaction()
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	player, err := trn.PlayerRead(cred.Username)
	trn.Abort()

	if err != nil {
		httpCode := http.StatusInternalServerError
		if _, ok := err.(store.DontExistError); ok {
			httpCode = http.StatusUnauthorized
		}

		w.WriteHeader(httpCode)

		return
	}

	if cred.VerifyHash([]byte(player.Password)) {
		session, err := s.sessionStore.Get(r, ServerName)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

			return
		}

		session.Values["username"] = player.Name
		session.Values["authenticated"] = true

		if err := session.Save(r, w); err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

			return
		}

		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func (s *Server) logoutUserHandler(w http.ResponseWriter, r *http.Request) {
	session, err := s.sessionStore.Get(r, ServerName)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	session.Values["authenticated"] = false

	if err := session.Save(r, w); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

		return
	}

	w.WriteHeader(http.StatusCreated)
}
