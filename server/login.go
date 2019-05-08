package server

import (
	"encoding/json"
	"net/http"
	"plateau/model"
	"plateau/server/response"
	"plateau/server/response/body"

	"github.com/jinzhu/gorm"
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
			response.WriteJSON(w, http.StatusInternalServerError, body.New().Ko(err))

			return
		}

		if session.Values["authenticated"] == nil || !session.Values["authenticated"].(bool) {
			w.WriteHeader(http.StatusForbidden)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) registerUser(w http.ResponseWriter, r *http.Request) {
	var cred loginCredentials

	if json.NewDecoder(r.Body).Decode(&cred) != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if !cred.IsValid() {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	resBody := body.New()

	hPassword, err := cred.PasswordHash()
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, resBody.Ko(err))

		return
	}

	if errs := s.db.Create(&model.Player{Name: cred.Username, Password: string(hPassword)}).GetErrors(); len(errs) > 0 {
		for _, err := range errs {
			if model.IsDuplicateError(err) {
				w.WriteHeader(http.StatusConflict)

				return
			}
		}

		response.WriteJSON(w, http.StatusInternalServerError, resBody.Ko(errs...))

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	var cred loginCredentials

	if json.NewDecoder(r.Body).Decode(&cred) != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if !cred.IsValid() {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var (
		resBody = body.New()

		player model.Player
	)

	if errs := s.db.Where("name = ?", cred.Username).Find(&player).GetErrors(); len(errs) > 0 {
		httpCode := http.StatusInternalServerError

		for _, err := range errs {
			if gorm.IsRecordNotFoundError(err) {
				httpCode = http.StatusUnauthorized

				break
			}
		}

		response.WriteJSON(w, httpCode, body.New().Ko(errs...))

		return
	}

	if cred.VerifyHash([]byte(player.Password)) {
		session, err := s.sessionStore.Get(r, ServerName)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, resBody.Ko(err))

			return
		}

		session.Values["username"] = player.Name
		session.Values["authenticated"] = true

		if err := session.Save(r, w); err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, resBody.Ko(err))

			return
		}

		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func (s *Server) logoutUser(w http.ResponseWriter, r *http.Request) {
	var (
		resBody = body.New()

		session, err = s.sessionStore.Get(r, ServerName)
	)

	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, resBody.Ko(err))

		return
	}

	session.Values["authenticated"] = false

	if err := session.Save(r, w); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, resBody.Ko(err))

		return
	}

	w.WriteHeader(http.StatusCreated)
}
