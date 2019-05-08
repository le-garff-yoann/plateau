package cmd

import (
	"log"
	"os"
	"os/signal"
	"plateau/server"
	"syscall"

	"github.com/spf13/cobra"
	// TODO: Add a little bit of "viper".
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Starts the server",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Starting the server....")

			srv, err := server.New(
				dbURL,
				serverListener,
				serverListenerSessionKey,
				serverListenerStaticDir,
				pgAutoMigrate,
				pgDebugging,
			)
			if err != nil {
				log.Fatal(err)
			}

			go func() {
				log.Fatal(srv.Start().Error())
			}()

			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

			log.Println("The server is started.")

			<-sigs

			log.Println("Gracefully stopping everything....")

			srv.Stop()

			os.Exit(0)
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.
		Flags().
		StringVarP(&dbURL, "pg-conn-str", "", "", "PostgreSQL connection string")
	runCmd.MarkFlagRequired("pg-conn-str")
	runCmd.
		Flags().
		BoolVarP(&pgAutoMigrate, "pg-automigrate", "", false, "enable PostgreSQL AutoMigration() (shouldn't be used in production)")
	runCmd.
		Flags().
		BoolVarP(&pgDebugging, "pg-debug", "", false, "enable PostgreSQL debugging")

	runCmd.
		Flags().
		StringVarP(&serverListener, "listen", "", ":8080", "listen on x:x (e.g. :8080 or 127.0.0.1:8080)")
	runCmd.
		Flags().
		StringVarP(&serverListenerSessionKey, "listen-session-key", "", "", "session key")
	runCmd.MarkFlagRequired("listen-session-key")
	// TODO: Add a switch to configure the session expiration (MaxAge).
	runCmd.
		Flags().
		StringVarP(&serverListenerStaticDir, "listen-static-dir", "", "", "exposes the contents of this directory at /")

	// TODO: Replace "log" with "github.com/apsdehal/go-logger".
	log.SetOutput(os.Stdout)
}
