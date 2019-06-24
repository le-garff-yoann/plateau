package cmd

import (
	"os"
	"os/signal"
	"plateau/server"
	"syscall"

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

			srv, err := server.New(
				serverListener, serverListenerStaticDir,
				gm,
				str,
			)
			if err != nil {
				logrus.Fatal(err)
			}

			go func() {
				logrus.Fatal(srv.Start().Error())
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
		StringVarP(&serverListenerStaticDir, "listen-static-dir", "", "", "Exposes the contents of this directory at /")

	runCmd.Flags().Var(&logLevel, "log-level", "Logrus log level")

	gm = newGame()
	if err := gm.Init(); err != nil {
		logrus.Fatal(err)
	}

	str = newStore()
	str.RunCommandSetter(runCmd)
}
