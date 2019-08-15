package cmd

import (
	"os"
	"os/signal"
	"plateau/server"
	"syscall"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	// TODO: Add a little bit of "viper".
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Starts the server",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.SetLevel(logrus.Level(logLevel))

			logrus.Info("Starting the server....")

			var byteSessionKeys [][]byte
			for _, sessionKey := range sessionKeys {
				byteSessionKeys = append(byteSessionKeys, []byte(sessionKey))
			}

			sessionStore := sessions.NewCookieStore(byteSessionKeys...)
			sessionStore.MaxAge(sessionMaxAge)

			srv, err := server.Init(
				serverListener, serverListenerStaticDir,
				gm,
				str,
				sessionStore,
			)
			if err != nil {
				logrus.Fatal(err)
			}

			if err := srv.Start(); err != nil {
				logrus.Fatal(err)
			}

			go func() {
				logrus.Fatal(srv.Listen().Error())
			}()

			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

			logrus.Info("The server is started.")

			<-sigs

			logrus.Info("Gracefully stopping everything....")

			if err := srv.Stop(); err != nil {
				logrus.Fatal(err)
			}

			os.Exit(0)
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.
		Flags().
		StringVarP(&serverListener, "listen", "l", ":8080", "Listen on x:x (e.g. :8080 or 127.0.0.1:8080)")
	runCmd.
		Flags().
		StringVar(&serverListenerStaticDir, "listen-static-dir", serverListenerStaticDir, "Exposes the contents of this directory at /")

	runCmd.
		Flags().
		StringArrayVar(&sessionKeys, "session-key", sessionKeys, `Session ("secret") key`)
	runCmd.MarkFlagRequired("session-key")
	runCmd.
		Flags().
		IntVar(&sessionMaxAge, "session-max-age", sessionMaxAge, "Sets the maximum duration of cookies in seconds")

	runCmd.Flags().Var(&logLevel, "log-level", "Logrus log level")

	gm = newGame()

	str = newStore()
	str.RunCommandSetter(runCmd)
}
