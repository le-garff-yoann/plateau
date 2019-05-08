package cmd

import (
	"fmt"
	"os"
	"plateau/server"

	"github.com/spf13/cobra"
)

// AppName is the CLI app name.
const AppName = server.ServerName

// RootCmd is meant to reused across cmd/*/*.go
var (
	dbURL, serverListener, serverListenerSessionKey, serverListenerStaticDir string
	pgDebugging, pgAutoMigrate                                               bool

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
