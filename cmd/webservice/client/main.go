package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

func main() {
	rootContext := context.Background()
	procContext, cancelSignalListeners := signal.NotifyContext(rootContext, os.Interrupt, unix.SIGTERM)
	defer cancelSignalListeners()

	if err := safeSession(procContext, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return
	}
}

func safeSession(ctx context.Context, log io.Writer) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(log, "panic: %s", err)
		}
	}()

	client := newWebClient()
	scenarios, err := client.GetScenarios(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", scenarios)
	scenarioName := scenarios.Available[0].Name
	//scenarioName := "multiple-up-and-back"

	controllers, err := client.GetControllers(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", controllers)

	reply, err := client.PostSession(ctx, PostSessionRequestBody{
		Scenario:   &scenarioName,
		Controller: &controllers.Available[0].Name,
	})
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", reply)

	for {
		reply, err := client.PostSessionTick(ctx, reply.SessionID)
		if err != nil {
			return err
		}
		if reply.Completed {
			break
		}
	}

	// dump log
	eventsReply, err := client.GetSessionEvents(ctx, reply.SessionID)
	if err != nil {
		return err
	}
	for _, e := range eventsReply.Events {
		switch e.EventType {
		case "InitStart":
			fmt.Printf("\t*** Init Start\n")
		case "InitDone":
			fmt.Printf("\t*** Init done.\n")
		case "InformFloor":
			fmt.Printf("\tFloor %d\n", *e.Floor)
		case "InformElevator":
			fmt.Printf("\tElevator %d\n", *e.Elevator)
		case "TickStart":
			fmt.Printf("\tTick Start %d\n", *e.Timestamp)
		case "ElevatorCalled":
			fmt.Printf("\t\tElevator called on floor %d\n", *e.Floor)
		case "TickDone":
			fmt.Printf("\tdone (%d)\n", *e.Timestamp)
		case "ElevatorFloorRequest":
			fmt.Printf("\t\tActor in elevator %d requesting floor %d\n", *e.Elevator, *e.Floor)
		case "ActorFinished":
			fmt.Printf("\t\tActor finshed.  Earned %d point(s).\n", *e.Points)
		case "ElevatorArrived":
			fmt.Printf("\t\tElevator %d arrived at floor %d\n", *e.Elevator, *e.Floor)
		case "ElevatorAtFloor":
			fmt.Printf("\t\tElevator %d is at floor %d\n", *e.Elevator, *e.Floor)
		default:
			fmt.Printf("\t\tUnhandled event: %+v\n", e)
		}
	}
	return nil
}
