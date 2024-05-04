package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tp-link",
	Short: "A CLI Client for TP-Link Router",
	Long: `Currently implements only SMS Inbox list and send commands.
	Tested only on Archer MR600`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
