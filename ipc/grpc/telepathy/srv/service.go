package srv

import (
	"context"
	"errors"
	"fmt"
	"github.com/meschbach/elevatinator/ipc/grpc/telepathy/pb"
	"github.com/meschbach/elevatinator/simulator"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

const (
	rcInit = iota
	rcRunning
	rcReset
)

type remoteController struct {
	pb.UnimplementedControllerServiceServer
	state int8
	Builder simulator.ControllerFunc
	Timeout time.Duration

	controller *controllerInstance
}

func newRemoteController(builder simulator.ControllerFunc)  *remoteController {
	return &remoteController{
		state:   rcInit,
		Builder: builder,
		Timeout: time.Second * 30,
	}
}

func (t *remoteController) Spawn(ctx context.Context, opts *pb.SpawnOptions) (*pb.Controller, error) {
	controller := &controllerInstance{
		pending:    make([]*pendingMove, 0),
	}
	controller.controller = t.Builder(controller)
	t.controller = controller
	return &pb.Controller{Id: 0}, nil
}

func doInit(t *remoteController, msg *pb.SimulationEvent_Init) error {
	elevatorCount := msg.ElevatorCount
	t.controller.maxElevators = elevatorCount
	//fake elevator ids for now
	elevatorIDs := make([]simulator.ElevatorID,elevatorCount)
	for i, _ := range elevatorIDs {
		elevatorIDs[i] = simulator.ElevatorID(i)
	}
	//dispatch to client
	t.controller.controller.Init(elevatorIDs)
	return nil
}

func (t *remoteController) Notice(ctx context.Context, notice *pb.SimulationNotice) (*pb.ControllerUpdates, error) {
	id := notice.Target.Id
	if id != 0 { return nil, errors.New("bad id") }

	fmt.Printf("Events: %#v\n", notice.Event)
	for _, e := range notice.Event {
		fmt.Println("Event")
		if e.Initialize != nil {
			if err := doInit(t,e.Initialize); err != nil {
				return nil, err
			}
		}

		if e.FloorSelection != nil {
			fmt.Println("Floor selection")
			elevator := e.FloorSelection.InElevator.ElevatorIndex
			floor := e.FloorSelection.Selected.FloorIndex
			t.controller.controller.FloorSelected(simulator.ElevatorID(elevator), simulator.FloorID(floor))
		}
	}

	out := make([]*pb.ControllerDirective, len(t.controller.pending))
	for i, e := range t.controller.pending {
		elevator := uint32(e.which)
		floor := uint32(e.to)
		fmt.Printf("Move elevator %d to %d\n", elevator, floor)
		out[i] = &pb.ControllerDirective{
			SeekFloor: &pb.ControllerDirective_MoveTo{
				Which: &pb.Elevator{ElevatorIndex: elevator},
				Target: &pb.Floor{FloorIndex: floor},
			},
		}
	}
	return &pb.ControllerUpdates{Pending: out}, nil
}

func RunControllerService(builder simulator.ControllerFunc)  {
	port := ":9998"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterControllerServiceServer(s, newRemoteController(builder))
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type pendingMove struct {
	which simulator.ElevatorID
	to simulator.FloorID
}

type controllerInstance struct {
	controller simulator.Controller
	pending []*pendingMove
	maxElevators uint32
}

func (c *controllerInstance) MoveTo(elevator simulator.ElevatorID, floor simulator.FloorID) {
	if elevator < 0 || uint32(elevator) > c.maxElevators {
		//TODO: Report problem
		panic(fmt.Sprintf("no such elevator %d", elevator))
	}
	fmt.Printf("Queuing move of %d to %d\n",elevator,floor)
	c.pending = append(c.pending, &pendingMove{
		which: elevator,
		to:    floor,
	})
}