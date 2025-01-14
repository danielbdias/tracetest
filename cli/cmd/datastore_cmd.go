package cmd

import (
	"fmt"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/spf13/cobra"
)

var dataStoreCmd = &cobra.Command{
	Use:    "datastore",
	Short:  "Manage your tracetest data stores",
	Long:   "Manage your tracetest data stores",
	PreRun: setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Track("Datastore", "cmd", map[string]string{})

		fmt.Println("Manage your data stores")
	},
	PostRun: teardownCommand,
}

func init() {
	rootCmd.AddCommand(dataStoreCmd)
}
