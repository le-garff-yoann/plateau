package server

import (
	"net/http"
	"plateau/store"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// ServerName is the server name.
const ServerName = "plateau"

// Server is basically the *plateau* runtime.
type Server struct {
	game Game

	matchRuntimesMux sync.Mutex
	matchRuntimes    map[string]*matchRuntime

	store        store.Store
	sessionStore sessions.Store

	router     *mux.Router
	httpServer *http.Server

	doneWg sync.WaitGroup
}

// New initializes the `Game` *gm* and returns a new `Server`.
func New(gm Game, str store.Store) (*Server, error) {
	s := &Server{
		game:          gm,
		matchRuntimes: make(map[string]*matchRuntime),
		store:         str,
	}

	if err := gm.Init(); err != nil {
		return nil, err
	}

	return s, nil
}

// Init returns a new `Server` with in addition the *gorilla/mux* router and the *sessionStore* initialized.
func Init(listener string, gm Game, str store.Store, sessionStore ...sessions.Store) (*Server, error) {
	s, err := New(gm, str)
	if err != nil {
		return nil, err
	}

	if len(sessionStore) == 0 {
		sessionStore = append(sessionStore, sessions.NewCookieStore([]byte("")))
	}
	s.sessionStore = sessionStore[0]

	s.router = mux.NewRouter().StrictSlash(true)
	s.httpServer = &http.Server{
		Addr:    listener,
		Handler: s.router,
	}

	ar := s.router.PathPrefix("/api").Subrouter()

	ar.Use(s.loginMiddleware)
	ar.
		PathPrefix("/game").
		Methods("GET").
		HandlerFunc(s.getGameDefinitionHandler).
		Name("readGame")
	ar.
		PathPrefix("/players/{name}").
		Methods("GET").
		HandlerFunc(s.readPlayerHandler).
		Name("readPlayer")
	ar.
		PathPrefix("/players").
		Methods("GET").
		HandlerFunc(s.getPlayersNameHandler).
		Name("getPlayersName")
	ar.
		PathPrefix("/matchs/{id}/players").
		Methods("GET").
		HandlerFunc(s.getMatchPlayersNameHandler).
		Name("getMatchPlayersName")
	ar.
		PathPrefix("/matchs/{id}/notifications").
		Methods("GET").
		HandlerFunc(s.streamMatchNotificationsHandler).
		Name("streamMatchNotifications")
	ar.
		PathPrefix("/matchs/{id}/deals").
		Methods("GET").
		HandlerFunc(s.getMatchDealsHandler).
		Name("getMatchDeals")
	ar.
		PathPrefix("/matchs/{id}").
		Methods("GET").
		HandlerFunc(s.readMatchHandler).
		Name("readMatch")
	ar.
		PathPrefix("/matchs/{id}").
		Methods("PATCH").
		HandlerFunc(s.patchMatchHandler).
		Name("patchMatch")
	ar.
		PathPrefix("/matchs").
		Methods("GET").
		HandlerFunc(s.getMatchIDsHandler).
		Name("getMatchIDs")
	ar.
		PathPrefix("/matchs").
		Methods("POST").
		HandlerFunc(s.createMatchHandler).
		Name("createMatch")

	s.router.
		PathPrefix("/user/register").
		Methods("POST").
		HandlerFunc(s.registerUserHandler).
		Name("registerUser")
	s.router.
		PathPrefix("/user/login").
		Methods("POST").
		HandlerFunc(s.loginUserHandler).
		Name("loginUser")
	s.router.
		PathPrefix("/user/logout").
		Methods("DELETE").
		HandlerFunc(s.logoutUserHandler).
		Name("logoutUser")

	return s, nil
}

// Start starts the server.
func (s *Server) Start() error {
	if err := s.store.Open(); err != nil {
		return err
	}

	s.doneWg.Add(1)

	return nil
}

// Stop stops the server.
func (s *Server) Stop() error {
	s.doneWg.Done()

	s.doneWg.Wait()

	return s.store.Close()
}

// Listen starts the HTTP server.
func (s *Server) Listen() error {
	return s.httpServer.ListenAndServe()
}
