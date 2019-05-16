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
	// serverStore *store.Store

	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Starts the server",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Starting the server....")

			srv, err := server.New(
				serverListener, serverListenerStaticDir,
				str,
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
		StringVarP(&serverListener, "listen", "l", ":8080", "Listen on x:x (e.g. :8080 or 127.0.0.1:8080)")
	runCmd.
		Flags().
		StringVarP(&serverListenerStaticDir, "listen-static-dir", "", "", "Exposes the contents of this directory at /")

	// TODO: Replace "log" with "github.com/apsdehal/go-logger".
	log.SetOutput(os.Stdout)

	str = newStore()
	str.RunCommandSetter(runCmd)
}
