package rethinkdb

import (
	"plateau/store"

	"github.com/gorilla/sessions"
	"github.com/le-garff-yoann/rethinkstore"
	"github.com/spf13/cobra"
	"github.com/stretchr/stew/slice"
	rethinkdb "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

var (
	address, database, sessionKey, username, password, authkey string
	maxIdle                                                    = 1
	maxOpen                                                    = 1

	automigrate = false
)

// Store ...
type Store struct {
	queryExecuter rethinkdb.QueryExecutor

	playerStore  *PlayerStore
	matchStore   *MatchStore
	sessionStore sessions.Store
}

// Open ...
func (s *Store) Open() error {
	var (
		queryExecuter rethinkdb.QueryExecutor

		err error
	)

	queryExecuter, err = rethinkdb.Connect(rethinkdb.ConnectOpts{
		Address: address, Database: database,
		Username: username, Password: password,
		AuthKey: authkey,
		MaxIdle: maxIdle, MaxOpen: maxOpen,
	})
	if err != nil {
		return err
	}

	s.queryExecuter = queryExecuter
	s.playerStore = &PlayerStore{s.queryExecuter}
	s.matchStore = &MatchStore{s.queryExecuter}

	s.sessionStore, err = rethinkstore.NewRethinkStore(
		address, database,
		"gorilla_sessions",
		maxIdle, maxOpen,
		[]byte(sessionKey),
	)
	if err != nil {
		return err
	}

	if automigrate {
		cursor, err := rethinkdb.TableList().Run(s.queryExecuter)
		if err != nil {
			return err
		}

		var existingTables []string
		if err := cursor.All(&existingTables); err != nil {
			return err
		}

		for _, tableName := range []string{
			s.playerStore.tableName(), s.matchStore.tableName(),
		} {
			if !slice.ContainsString(existingTables, tableName) {
				if err := rethinkdb.TableCreate(tableName).Exec(s.queryExecuter); err != nil {
					return err
				}
			}
		}
	}

	return err
}

// Close ...
func (s *Store) Close() error {
	var err error

	session, ok := s.queryExecuter.(*rethinkdb.Session)
	if ok {
		err = session.Close()
	}

	sessionStore, ok := s.sessionStore.(*rethinkstore.RethinkStore)
	if ok {
		sessionStore.Close()
	}

	return err
}

// RunCommandSetter ...
func (s *Store) RunCommandSetter(runCmd *cobra.Command) {
	runCmd.
		Flags().
		StringVarP(&address, "rethinkdb-address", "", address, "RethinkDB server address")
	runCmd.MarkFlagRequired("rethinkdb-address")
	runCmd.
		Flags().
		StringVarP(&database, "rethinkdb-database", "", database, "RethinkDB database name")
	runCmd.MarkFlagRequired("rethinkdb-database")
	runCmd.
		Flags().
		StringVarP(&database, "rethinkdb-username", "", username, "RethinkDB database username")
	runCmd.
		Flags().
		StringVarP(&database, "rethinkdb-password", "", password, "RethinkDB database password")
	runCmd.
		Flags().
		IntVarP(&maxIdle, "rethinkdb-max-idle", "", maxIdle, "RethinkDB max idle connctions")
	runCmd.
		Flags().
		IntVarP(&maxIdle, "rethinkdb-max-open", "", maxOpen, "RethinkDB max open connections")
	runCmd.
		Flags().
		BoolVarP(&automigrate, "rethinkdb-create-tables", "", automigrate, "Create tables at startup")

	runCmd.
		Flags().
		StringVarP(&sessionKey, "session-key", "", sessionKey, `Session ("secret") key`)
	runCmd.MarkFlagRequired("session-key")
	// TODO: Add a switch to configure the session expiration (MaxAge).
}

// Players ...
func (s *Store) Players() store.PlayerStore {
	return s.playerStore
}

// Matchs ...
func (s *Store) Matchs() store.MatchStore {
	return s.matchStore
}

// Sessions ...
func (s *Store) Sessions() sessions.Store {
	return s.sessionStore
}
