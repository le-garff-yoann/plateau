package server

import (
	"net/http"
	"plateau/broadcaster"
	"plateau/store"
	"sync"

	"github.com/gorilla/mux"
)

// ServerName is the server name.
const ServerName = "plateau"

// Server ...
type Server struct {
	game Game

	matchmatchRuntimesMux sync.Mutex
	matchRuntimes         map[string]*matchRuntime

	store store.Store

	router     *mux.Router
	httpServer *http.Server

	doneBroadcaster *broadcaster.Broadcaster
	doneWg          sync.WaitGroup
}

// New ...
func New(listener, listenerStaticDir string, gm Game, str store.Store) (*Server, error) {
	if err := gm.Init(); err != nil {
		return nil, err
	}

	if err := str.Open(); err != nil {
		return nil, err
	}

	r := mux.NewRouter().StrictSlash(true)
	ar := r.PathPrefix("/api").Subrouter()

	s := &Server{
		game:          gm,
		matchRuntimes: make(map[string]*matchRuntime),
		store:         str,
		router:        r,
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
		PathPrefix("/matchs/{id}/players").
		HandlerFunc(s.getMatchPlayersNameHandler).
		Name("getMatchPlayersName")
	ar.
		PathPrefix("/matchs/{id}/deals").
		HandlerFunc(s.getMatchDealsHandler).
		Name("getMatchDeals")
	ar.
		PathPrefix("/matchs/{id}/notifications").
		HandlerFunc(s.streamMatchNotificationsHandler).
		Name("streamMatchNotifications")
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

	return nil
}

// Stop ...
func (s *Server) Stop() error {
	s.doneBroadcaster.Submit(0)
	s.doneWg.Wait()

	s.doneBroadcaster.Done()

	return s.store.Close()
}

// Listen ...
func (s *Server) Listen() error {
	return s.httpServer.ListenAndServe()
}
