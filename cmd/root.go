package cmd

import (
	"fmt"
	"os"
	"plateau/pflag"
	"plateau/server"
	"plateau/store"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AppName is the CLI app name.
const AppName = server.ServerName

// RootCmd is meant to reused across cmd/*/*.go
var (
	serverListener, serverListenerStaticDir string
	logLevel                                = pflag.LogrusLevel(logrus.InfoLevel)
	gm                                      server.Game
	str                                     store.Store

	rootCmd = &cobra.Command{
		Use: AppName,
	}
)

// Execute execute `rootCmd`.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}
