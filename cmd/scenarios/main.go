package main

import (
	"fmt"
	"github.com/meschbach/elevatinator/pkg/ipc/grpc/telepathy"
	"github.com/meschbach/elevatinator/pkg/scenarios"
	"github.com/meschbach/elevatinator/pkg/simulator"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	serviceAddress := "localhost:9998"

	runScenario := func(use string, short string, setup func(simulation *simulator.Simulation) simulator.Tick) *cobra.Command {
		return &cobra.Command{
			Use:   use,
			Short: short,
			RunE: func(cmd *cobra.Command, args []string) error {
				bridge, err := telepathy.DialLanding(serviceAddress)
				if err != nil {
					return err
				}

				scenarios.RunScenario(bridge.ControllerAdapter(), setup)
				return nil
			},
		}
	}

	rootCmd := &cobra.Command{
		Use:   "scenarios",
		Short: "Run scenarios against an AI gRPC service",
	}
	rootCmd.PersistentFlags().StringVarP(&serviceAddress, "ai-address", "a", serviceAddress, "AI unit address to connect to")
	rootCmd.AddCommand(runScenario("single-up", "Runs a scenario for a single person to go up", scenarios.SinglePersonUp))
	rootCmd.AddCommand(runScenario("single-down", "Runs a scenario for a single person to go down", scenarios.SinglePersonDown))
	rootCmd.AddCommand(runScenario("multiple-up-and-back", "Runs a scenario with various persons going up and back", scenarios.MultipleUpAndBack))
	rootCmd.AddCommand(healthProbeCommand(&serviceAddress))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func healthProbeCommand(serviceAddress *string) *cobra.Command {
	return &cobra.Command{
		Use:   "health-probe",
		Short: "Uses standard gRPC health check to ensure a service is healthy",
		RunE: func(cmd *cobra.Command, args []string) error {
			address := *serviceAddress
			//todo: add implicit retry
			if healthy, err := telepathy.CheckHealth(cmd.Context(), address); err == nil {
				if healthy {
					fmt.Printf("%q is healthy\n", address)
				} else {
					fmt.Printf("UNHEALTHY:\t\t%q\n", address)
				}
				return nil
			} else {
				fmt.Printf("Unable to check health:\t\t%q\t%e\n", address, err)
				return err
			}
		},
	}
}
