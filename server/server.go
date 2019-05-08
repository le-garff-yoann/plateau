package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"plateau/model"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/wader/gormstore"
)

// ServerName is the server name.
const ServerName = "plateau"

var wsUpgrader = websocket.Upgrader{}

// Server ...
type Server struct {
	db *gorm.DB
	// store        *store.Store
	sessionStore *gormstore.Store

	HTTPServer *http.Server
}

// New ...
func New(pgURL, listener, listenerSessionKey, listenerStaticDir string, pgAutoMigrate, pgDebugging bool) (*Server, error) {
	db, err := gorm.Open("postgres", pgURL)
	if err != nil {
		return nil, err
	}

	db.LogMode(pgDebugging)
	if pgDebugging {
		db.SetLogger(log.New(os.Stdout, "\r\n", 0))
	}

	if pgAutoMigrate {
		if errs := model.AutoMigrate(db).GetErrors(); len(errs) > 0 {
			return nil, fmt.Errorf("%s", errs)
		}
	}

	r := mux.NewRouter().StrictSlash(true)
	ar := r.PathPrefix("/api").Subrouter()

	s := &Server{
		db:           db,
		sessionStore: gormstore.New(db, []byte(listenerSessionKey)),
		HTTPServer: &http.Server{
			Addr:    listener,
			Handler: r,
		},
	}

	ar.Use(s.loginMiddleware)
	ar.
		PathPrefix("/players/{name}/games").
		Methods("GET").
		HandlerFunc(s.getPlayerGamesID).
		Name("getPlayerGamesID")
	ar.
		PathPrefix("/players/{name}").
		Methods("GET").
		HandlerFunc(s.readPlayer).
		Name("readPlayer")
	ar.
		PathPrefix("/players").
		Methods("GET").
		HandlerFunc(s.getPlayersName).
		Name("getPlayersName")
	ar.
		PathPrefix("/games/{id}/players").
		HandlerFunc(s.getGamePlayersName).
		Name("getGamePlayersName")
	ar.
		PathPrefix("/games/{id}/event-containers").
		HandlerFunc(s.getGameEventContainers).
		Name("getGameEventContainers")
	ar.
		PathPrefix("/games/{id}").
		Headers("X-Interactive", "true").
		HandlerFunc(s.connectGame).
		Name("connectGame")
	ar.
		PathPrefix("/games/{id}").
		Methods("GET").
		HandlerFunc(s.readGame).
		Name("readGame")
	ar.
		PathPrefix("/games").
		Methods("GET").
		HandlerFunc(s.getGameIDs).
		Name("getGameIDs")
	ar.
		PathPrefix("/games").
		Methods("POST").
		HandlerFunc(s.createGame).
		Name("createGame")

	r.
		PathPrefix("/user/register").
		Methods("POST").
		HandlerFunc(s.registerUser).
		Name("registerUser")
	r.
		PathPrefix("/user/login").
		Methods("POST").
		HandlerFunc(s.loginUser).
		Name("loginUser")
	r.
		PathPrefix("/user/logout").
		Methods("DEL").
		HandlerFunc(s.logoutUser).
		Name("logoutUser")
	// TODO: Add an endpoint from which we can get the game name and it's default settings.

	if listenerStaticDir != "" {
		r.
			PathPrefix("/").
			Methods("GET").
			Handler(http.FileServer(http.Dir(listenerStaticDir))).
			Name("root")
	}

	return s, nil
}

// Start ...
func (s *Server) Start() error {
	return s.HTTPServer.ListenAndServe()
}

// Stop ...
func (s *Server) Stop() error {
	return s.db.Close()
}
