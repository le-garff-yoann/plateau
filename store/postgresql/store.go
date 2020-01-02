package postgresql

import (
	"encoding/json"
	"fmt"
	"plateau/broadcaster"
	"plateau/server"
	"plateau/store"

	"github.com/sirupsen/logrus"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/spf13/cobra"
)

var (
	postgresqlURL          = ""
	createTableOptionsTemp = false

	pubSubChName = server.ServerName
)

// Store implements the `store.Store` interface.
type Store struct {
	db *pg.DB
	ln *pg.Listener

	matchNotificationsBroadcaster *broadcaster.Broadcaster

	done chan int
}

// Open implements the `store.Store` interface.
func (s *Store) Open() (err error) {
	pgOptions, err := pg.ParseURL(postgresqlURL)
	if err != nil {
		return err
	}

	s.db = pg.Connect(pgOptions)

	s.done = make(chan int)

	s.matchNotificationsBroadcaster = broadcaster.New()
	go s.matchNotificationsBroadcaster.Run()

	s.ln = s.db.Listen(pubSubChName)
	notifCh := s.ln.Channel()

	go func() {
		for {
			select {
			case r := <-notifCh:
				var n store.MatchNotification
				if err := json.Unmarshal([]byte(r.Payload), &n); err == nil {
					s.matchNotificationsBroadcaster.Submit(n)
				} else {
					logrus.Error(err)
				}
			case <-s.done:
				return
			}
		}
	}()

	// REFACTOR: Need a production-ready schema migration method.
	for _, model := range []interface{}{
		&Match{}, &Player{}, &MatchPlayer{},
	} {
		if err := s.db.CreateTable(model, &orm.CreateTableOptions{
			Temp:          createTableOptionsTemp,
			IfNotExists:   true,
			FKConstraints: true,
		}); err != nil {
			return err
		}
	}

	return nil
}

// Close implements the `store.Store` interface.
func (s *Store) Close() error {
	s.done <- 0
	s.matchNotificationsBroadcaster.Done()

	if err := s.ln.Close(); err != nil {
		return err
	}

	return s.db.Close()
}

// RunCommandSetter implements the `store.Store` interface.
func (s *Store) RunCommandSetter(runCmd *cobra.Command) {
	runCmd.
		Flags().
		StringVar(&postgresqlURL, "pg-url", postgresqlURL, "PostgreSQL URL")
	runCmd.MarkFlagRequired("pg-url")
}

// BeginTransaction implements the `store.Store` interface.
func (s *Store) BeginTransaction() (store.Transaction, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	return &Transaction{
		tx: tx,
		commitCb: func(trn *Transaction) error {
			for _, matchNotification := range trn.matchNotifications {
				json, err := json.Marshal(matchNotification)
				if err != nil {
					logrus.Fatal(err)
				}

				if _, err := s.db.Exec(fmt.Sprintf("NOTIFY %s, ?", pubSubChName), string(json)); err != nil {
					logrus.Warn(err)
				}
			}

			return nil
		},
		closed: false,
	}, nil
}

// RegisterNotificationsChannel implements the `store.Store` interface.
func (s *Store) RegisterNotificationsChannel(ch chan interface{}) error {
	s.matchNotificationsBroadcaster.Register(ch)

	return nil
}

// UnregisterNotificationsChannel implements the `store.Store` interface.
func (s *Store) UnregisterNotificationsChannel(ch chan interface{}) error {
	s.matchNotificationsBroadcaster.Unregister(ch)

	return nil
}
