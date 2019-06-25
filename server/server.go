package server

import (
	"net/http"
	"plateau/broadcaster"
	"plateau/store"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// ServerName is the server name.
const ServerName = "plateau"

var wsUpgrader = websocket.Upgrader{}

// Server ...
type Server struct {
	game Game

	matchmatchRuntimesMux sync.Mutex
	matchRuntimes         map[string]*MatchRuntime

	store store.Store

	httpServer *http.Server

	doneBroadcaster *broadcaster.Broadcaster
	doneWg          sync.WaitGroup
}

// New ...
func New(listener, listenerStaticDir string, gm Game, str store.Store) (*Server, error) {
	if err := str.Open(); err != nil {
		return nil, err
	}

	r := mux.NewRouter().StrictSlash(true)
	ar := r.PathPrefix("/api").Subrouter()

	s := &Server{
		game:          gm,
		matchRuntimes: make(map[string]*MatchRuntime),
		store:         str,
		httpServer: &http.Server{
			Addr:    listener,
			Handler: r,
		},
		doneBroadcaster: broadcaster.New(),
	}

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
		PathPrefix("/matchs/{id}/connected-players").
		HandlerFunc(s.getMatchConnectedPlayersNameHandler).
		Name("getMatchPlayersName")
	ar.
		PathPrefix("/matchs/{id}/players").
		HandlerFunc(s.getMatchPlayersNameHandler).
		Name("getMatchPlayersName")
	ar.
		PathPrefix("/matchs/{id}/deals").
		HandlerFunc(s.getMatchDealsHandler).
		Name("getMatchEventContainers")
	ar.
		PathPrefix("/matchs/{id}").
		Headers("X-Interactive", "true").
		HandlerFunc(s.connectMatchHandler).
		Name("connectMatch")
	ar.
		PathPrefix("/matchs/{id}").
		Methods("GET").
		HandlerFunc(s.readMatchHandler).
		Name("readMatch")
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

	r.
		PathPrefix("/user/register").
		Methods("POST").
		HandlerFunc(s.registerUserHandler).
		Name("registerUser")
	r.
		PathPrefix("/user/login").
		Methods("POST").
		HandlerFunc(s.loginUserHandler).
		Name("loginUser")
	r.
		PathPrefix("/user/logout").
		Methods("DEL").
		HandlerFunc(s.logoutUserHandler).
		Name("logoutUser")

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
	go s.doneBroadcaster.Run()

	s.doneWg.Add(1)
	ch, uuid := s.doneBroadcaster.Subscribe()

	go func() {
		<-ch

		s.doneBroadcaster.Unsubscribe(uuid)
		s.doneWg.Done()
	}()

	return s.httpServer.ListenAndServe()
}

// Stop ...
func (s *Server) Stop() error {
	s.doneBroadcaster.Submit(0)
	s.doneWg.Wait()

	s.doneBroadcaster.Done()

	return s.store.Close()
}
