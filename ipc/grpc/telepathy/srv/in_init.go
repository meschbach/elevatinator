package srv

import (
	"github.com/meschbach/elevatinator/ipc/grpc/telepathy/pb"
	"github.com/meschbach/elevatinator/simulator"
)

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
