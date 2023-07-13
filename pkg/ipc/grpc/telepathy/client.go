package telepathy

import (
	"context"
	"fmt"
	pb2 "github.com/meschbach/elevatinator/pkg/ipc/grpc/telepathy/pb"
	simulator2 "github.com/meschbach/elevatinator/pkg/simulator"
	"google.golang.org/grpc"
	"time"
)

type ConnectionError struct {
	Target     string
	Underlying error
}

func (c *ConnectionError) Unwrap() error {
	return c.Underlying
}

func (c *ConnectionError) Error() string {
	return fmt.Sprintf("failed to connect to %q", c.Target)
}

type Landing struct {
	connection *grpc.ClientConn
	client     pb2.ControllerServiceClient
}

func DialLanding(address string) (*Landing, error) {
	// Set up a connection to the server.
	fmt.Println("Attempting to connect")
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		return nil, &ConnectionError{
			Target:     address,
			Underlying: err,
		}
	}
	return LandingWithConnection(conn), nil
}

func LandingWithConnection(conn *grpc.ClientConn) *Landing {
	c := pb2.NewControllerServiceClient(conn)

	return &Landing{
		connection: conn,
		client:     c,
	}
}

func (l *Landing) ControllerAdapter() simulator2.ControllerFunc {
	return func(elevators simulator2.ControlledElevators) simulator2.Controller {
		ctx, done := context.WithTimeout(context.Background(), time.Second*1)
		defer done()

		result, err := l.client.Spawn(ctx, &pb2.SpawnOptions{})
		if err != nil {
			panic(err)
		}
		return &BridgedController{
			controllerID: result.Id,
			landing:      l,
			controls:     elevators,
		}
	}
}

type BridgedController struct {
	landing      *Landing
	controllerID uint32
	controls     simulator2.ControlledElevators
	elevators    []simulator2.ElevatorID
}

func (m *BridgedController) Init(elevators []simulator2.ElevatorID) {
	m.elevators = elevators
	m.dispatch(&pb2.SimulationEvent{
		When: nil,
		Initialize: &pb2.SimulationEvent_Init{
			ElevatorCount: uint32(len(elevators)),
			FloorCount:    5,
		},
	})
}

func (m *BridgedController) Called(floor simulator2.FloorID) {
	m.dispatch(&pb2.SimulationEvent{
		Called: &pb2.SimulationEvent_ElevatorCalled{CalledAt: &pb2.Floor{FloorIndex: uint32(floor)}},
	})
}

func (m *BridgedController) FloorSelected(elevatorID simulator2.ElevatorID, floor simulator2.FloorID) {
	m.dispatch(&pb2.SimulationEvent{
		When: &pb2.Tick{V0: 0},
		FloorSelection: &pb2.SimulationEvent_FloorSelected{
			InElevator: &pb2.Elevator{ElevatorIndex: uint32(elevatorID)},
			Selected:   &pb2.Floor{FloorIndex: uint32(floor)},
		},
	})
}

func (m *BridgedController) CompletedMove(elevatorID simulator2.ElevatorID) {
	m.dispatch(&pb2.SimulationEvent{
		Arriving: &pb2.SimulationEvent_ElevatorArrived{
			Arriving: &pb2.Elevator{ElevatorIndex: uint32(elevatorID)},
		},
	})
}

func (m *BridgedController) dispatch(e *pb2.SimulationEvent) {
	ctx, done := context.WithTimeout(context.Background(), time.Second*1)
	defer done()

	updates, err := m.landing.client.Notice(ctx, &pb2.SimulationNotice{
		Target: &pb2.Controller{Id: m.controllerID},
		Event:  []*pb2.SimulationEvent{e},
	})
	if err != nil {
		panic(err)
	}

	for _, p := range updates.Pending {
		if p.SeekFloor != nil {
			floor := convertFloorFromWire(p.SeekFloor.Target)
			elevator, err := m.convertElevatorFromWire(p.SeekFloor.Which)
			if err != nil {
				//TODO: Should end the simulation
				panic(err)
			}
			fmt.Printf("Sending to floor %d for elevator %d\n", floor, elevator)

			m.controls.MoveTo(elevator, floor)
		}
	}
}

func convertFloorFromWire(input *pb2.Floor) simulator2.FloorID {
	index := input.FloorIndex
	fmt.Printf("Floor %d\n", index)
	return simulator2.FloorID(index)
}

func (m *BridgedController) convertElevatorFromWire(input *pb2.Elevator) (simulator2.ElevatorID, error) {
	index := input.ElevatorIndex
	if len(m.elevators) < int(index) {
		return -1, fmt.Errorf("got elevator index %d, max %d", index, len(m.elevators))
	}
	fmt.Printf("Elevator @ index %d\n", index)
	return m.elevators[index], nil
}
