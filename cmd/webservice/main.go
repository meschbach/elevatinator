package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	runCommand := &cobra.Command{
		Use:   "run",
		Short: "Runs a webservice to control and maintain interactions",
		RunE: func(cmd *cobra.Command, args []string) error {
			runService(cmd.Context())
			return nil
		},
	}

	rootCmd := &cobra.Command{
		Use:   "webservice",
		Short: "Runs the webservice for online visualization and interaction",
	}
	rootCmd.AddCommand(runCommand)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
