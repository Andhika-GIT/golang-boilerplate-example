package commands

import (
	"log"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "user-service",
	Short: "User Service CLI",
}

func RegisterCommand() {
	RootCmd.AddCommand(ServerCmd)
	RootCmd.AddCommand(MigrateCmd)
}

func Execute() {
	RegisterCommand()
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf(err.Error())
	}
}
