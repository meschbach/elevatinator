package main

import (
	"fmt"
	"github.com/meschbach/elevatinator/ipc/grpc/telepathy"
	"github.com/meschbach/elevatinator/scenarios"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	serviceAddress := "localhost:9998"

	singleUp := &cobra.Command{
		Use:   "single-up",
		Short: "Runs single person up scenario",
		RunE: func(cmd *cobra.Command, args []string) error {
			bridge, err := telepathy.DialLanding("localhost:9998")
			if err != nil {
				return err
			}

			//scenarios := []scenarios.Scenario{scenarios.SinglePersonUp}
			scenario := scenarios.SinglePersonDown
			scenarios.RunScenario(bridge.ControllerAdapter(), scenario)
			return nil
		},
	}

	rootCmd := &cobra.Command{
		Use:   "scenarios",
		Short: "Run scenarios against an AI gRPC service",
	}
	rootCmd.PersistentFlags().StringVarP(&serviceAddress, "ai-address", "a", serviceAddress, "AI unit address to connect to")
	rootCmd.AddCommand(singleUp)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
