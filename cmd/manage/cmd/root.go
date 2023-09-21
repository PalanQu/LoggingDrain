package cmd

import (
	"os"

	config_cmd "loggingdrain/cmd/manage/cmd/config"

	"github.com/spf13/cobra"
)

var dryrun bool

var rootCmd = &cobra.Command{
	Use:   "manage api",
	Short: "Manage the apis",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func InitCmd() error {
	rootCmd.PersistentFlags().BoolVarP(&dryrun, "dry-run", "n", false, "--dry-run")
	config_cmd.Init(rootCmd)
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
