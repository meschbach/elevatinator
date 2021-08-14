package telepathy

import (
	"context"
	"fmt"
	"github.com/meschbach/elevatinator/ipc/grpc/telepathy/pb"
	"github.com/meschbach/elevatinator/simulator"
	"google.golang.org/grpc"
	"time"
)

type Landing struct {
	connection *grpc.ClientConn
	client pb.ControllerServiceClient
}

func DialLanding(address string) (*Landing, error) {
	// Set up a connection to the server.
	fmt.Println("Attempting to connect")
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(1 * time.Second))
	if err != nil {
		return nil, err
	}
	c := pb.NewControllerServiceClient(conn)

	return &Landing{
		connection: conn,
		client:     c,
	}, nil
}

func (l *Landing) ControllerAdapter() simulator.ControllerFunc {
	return func(elevators simulator.ControlledElevators) simulator.Controller {
		ctx, done := context.WithTimeout(context.Background(), time.Second * 1)
		defer done()

		result, err := l.client.Spawn(ctx,&pb.SpawnOptions{})
		if err != nil { panic(err) }
		return &BridgedController{
			controllerID: result.Id,
			landing: l,
			controls: elevators,
		}
	}
}

type BridgedController struct {
	landing *Landing
	controllerID uint32
	controls simulator.ControlledElevators
	elevators []simulator.ElevatorID
}

func (m *BridgedController) Init(elevators []simulator.ElevatorID) {
	m.elevators = elevators
	m.dispatch(&pb.SimulationEvent{
		When:           nil,
		Initialize:     &pb.SimulationEvent_Init{
			ElevatorCount: uint32(len(elevators)),
			FloorCount:    5,
		},
	})
}

func (m *BridgedController) Called(floor simulator.FloorID) {
	m.dispatch(&pb.SimulationEvent{
		Called: &pb.SimulationEvent_ElevatorCalled{CalledAt: &pb.Floor{FloorIndex: uint32(floor)}},
	})
}

func (m *BridgedController) FloorSelected(elevatorID simulator.ElevatorID, floor simulator.FloorID) {
	m.dispatch(&pb.SimulationEvent{
		When:           &pb.Tick{V0: 0},
		FloorSelection: &pb.SimulationEvent_FloorSelected{
			InElevator: &pb.Elevator{ElevatorIndex: uint32(elevatorID)},
			Selected:   &pb.Floor{FloorIndex: uint32(floor)},
		},
	})
}

func (m *BridgedController) CompletedMove(elevatorID simulator.ElevatorID) {
	m.dispatch(&pb.SimulationEvent{
		Arriving:       &pb.SimulationEvent_ElevatorArrived{
			Arriving:   &pb.Elevator{ElevatorIndex: uint32(elevatorID)},
		},
	})
}

func (m *BridgedController) dispatch(e *pb.SimulationEvent) {
	ctx, done := context.WithTimeout(context.Background(), time.Second * 1)
	defer done()

	updates, err := m.landing.client.Notice(ctx, &pb.SimulationNotice{
		Target: &pb.Controller{Id: m.controllerID},
		Event:  []*pb.SimulationEvent{e},
	})
	if err != nil { panic(err) }

	for _, p := range updates.Pending {
		if p.SeekFloor != nil {
			floor := convertFloorFromWire(p.SeekFloor.Target)
			elevator, err := m.convertElevatorFromWire(p.SeekFloor.Which)
			if err != nil {
				//TODO: Should end the simulation
				panic(err)
			}
			fmt.Printf("Sending to floor %d for elevator %d\n", floor, elevator)

			m.controls.MoveTo(elevator,floor)
		}
	}
}

func convertFloorFromWire(input *pb.Floor) simulator.FloorID {
	index := input.FloorIndex
	fmt.Printf("Floor %d\n", index)
	return simulator.FloorID(index)
}

func (m *BridgedController) convertElevatorFromWire(input *pb.Elevator) (simulator.ElevatorID, error) {
	index := input.ElevatorIndex
	if len(m.elevators) < int(index) {
		return -1, fmt.Errorf("got elevator index %d, max %d", index, len(m.elevators))
	}
	fmt.Printf("Elevator @ index %d\n", index)
	return m.elevators[index],nil
}
