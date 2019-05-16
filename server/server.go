package server

import (
	"net/http"
	"plateau/store"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// ServerName is the server name.
const ServerName = "plateau"

var wsUpgrader = websocket.Upgrader{}

// Server ...
type Server struct {
	store store.Store

	ecBroadcastersMux sync.Mutex
	ecBroadcasters    map[string]*store.EventContainerBroadcaster

	httpServer *http.Server

	done chan int
}

func (s *Server) recvEventContainerBroadcaster(matchID string) (<-chan store.EventContainer, *uuid.UUID, error) {
	if _, ok := s.ecBroadcasters[matchID]; !ok {
		s.ecBroadcastersMux.Lock()
		defer s.ecBroadcastersMux.Unlock()

		var err error
		s.ecBroadcasters[matchID], err = s.store.Matchs().CreateEventContainerBroadcaster(matchID)
		if err != nil {
			delete(s.ecBroadcasters, matchID)

			return nil, nil, err
		}

		go func() {
			s.ecBroadcasters[matchID].Run()
		}()

		recv, uuid := s.ecBroadcasters[matchID].Recv()

		return recv, &uuid, nil
	}

	return nil, nil, nil
}

func (s *Server) removeRecvEventContainerBroadcaster(matchID string, uuid uuid.UUID) {
	s.ecBroadcastersMux.Lock()

	s.ecBroadcasters[matchID].RemoveRecv(uuid)

	s.ecBroadcastersMux.Unlock()
}

func (s *Server) cleanupBroadcasters() {
	s.ecBroadcastersMux.Lock()

	for matchID, br := range s.ecBroadcasters {
		if br.CountReceivers() == 0 {
			br.Done <- 0

			delete(s.ecBroadcasters, matchID)
		}
	}

	s.ecBroadcastersMux.Unlock()
}

// New ...
func New(listener, listenerStaticDir string, str store.Store) (*Server, error) {
	if err := str.Open(); err != nil {
		return nil, err
	}

	r := mux.NewRouter().StrictSlash(true)
	ar := r.PathPrefix("/api").Subrouter()

	s := &Server{
		store:          str,
		ecBroadcasters: make(map[string]*store.EventContainerBroadcaster),
		httpServer: &http.Server{
			Addr:    listener,
			Handler: r,
		},
	}

	ar.Use(s.loginMiddleware)
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
		PathPrefix("/matchs/{id}/event-containers").
		HandlerFunc(s.getMatchEventContainersHandler).
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
	// TODO: Add an endpoint from which we can get the game implementation name and it's default settings.

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
	tk := time.NewTicker(time.Second * 5)

	go func() {
		for {
			select {
			case <-tk.C:
				s.cleanupBroadcasters()
			case <-s.done:
				tk.Stop()

				return
			}
		}
	}()

	return s.httpServer.ListenAndServe()
}

// Stop ...
func (s *Server) Stop() error {
	s.done <- 0

	return s.store.Close()
}
