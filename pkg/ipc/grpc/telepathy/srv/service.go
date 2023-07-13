package srv

import (
	"context"
	"errors"
	"fmt"
	pb2 "github.com/meschbach/elevatinator/pkg/ipc/grpc/telepathy/pb"
	simulator2 "github.com/meschbach/elevatinator/pkg/simulator"
	"time"
)

const (
	rcInit = iota
	rcRunning
	rcReset
)

type remoteController struct {
	pb2.UnimplementedControllerServiceServer
	state   int8
	Builder simulator2.ControllerFunc
	Timeout time.Duration

	controller *controllerInstance
}

func newRemoteController(builder simulator2.ControllerFunc) *remoteController {
	return &remoteController{
		state:   rcInit,
		Builder: builder,
		Timeout: time.Second * 30,
	}
}

func (t *remoteController) Spawn(ctx context.Context, opts *pb2.SpawnOptions) (*pb2.Controller, error) {
	controller := &controllerInstance{
		pending: make([]*pendingMove, 0),
	}
	controller.controller = t.Builder(controller)
	t.controller = controller
	return &pb2.Controller{Id: 0}, nil
}

func (t *remoteController) Notice(ctx context.Context, notice *pb2.SimulationNotice) (*pb2.ControllerUpdates, error) {
	id := notice.Target.Id
	if id != 0 {
		return nil, errors.New("bad id")
	}

	fmt.Printf("Events: %#v\n", notice.Event)
	for _, e := range notice.Event {
		fmt.Println("Event")
		if e.Initialize != nil {
			if err := doInit(t, e.Initialize); err != nil {
				return nil, err
			}
		}

		if e.Called != nil {
			if err := doFloorCall(t, e.Called); err != nil {
				return nil, err
			}
		}
		if e.Arriving != nil {
			if err := doElevatorArrived(t, e.Arriving); err != nil {
				return nil, err
			}
		}

		if e.FloorSelection != nil {
			fmt.Println("Floor selection")
			elevator := e.FloorSelection.InElevator.ElevatorIndex
			floor := e.FloorSelection.Selected.FloorIndex
			t.controller.controller.FloorSelected(simulator2.ElevatorID(elevator), simulator2.FloorID(floor))
		}
	}

	out := make([]*pb2.ControllerDirective, len(t.controller.pending))
	for i, e := range t.controller.pending {
		elevator := uint32(e.which)
		floor := uint32(e.to)
		fmt.Printf("Move elevator %d to %d\n", elevator, floor)
		out[i] = &pb2.ControllerDirective{
			SeekFloor: &pb2.ControllerDirective_MoveTo{
				Which:  &pb2.Elevator{ElevatorIndex: elevator},
				Target: &pb2.Floor{FloorIndex: floor},
			},
		}
	}
	t.controller.resetPending()
	return &pb2.ControllerUpdates{Pending: out}, nil
}
