package main

import (
	"fmt"
	"github.com/meschbach/elevatinator/pkg/controllers/queue"
	"github.com/meschbach/elevatinator/pkg/ipc/grpc/telepathy/srv"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	serviceAddress := "localhost:9998"

	run := &cobra.Command{
		Use:   "run",
		Short: "launches the serivce",
		RunE: func(cmd *cobra.Command, args []string) error {
			return srv.RunControllerService(queue.NewController)
		},
	}

	rootCmd := &cobra.Command{
		Use:   "queue",
		Short: "Elevatinator AI unit using a queueing technique",
	}
	rootCmd.PersistentFlags().StringVarP(&serviceAddress, "address", "a", serviceAddress, "Binding address for runs")
	rootCmd.AddCommand(run)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
