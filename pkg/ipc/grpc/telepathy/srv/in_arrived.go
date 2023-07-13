package srv

import (
	"github.com/meschbach/elevatinator/pkg/ipc/grpc/telepathy/pb"
	"github.com/meschbach/elevatinator/pkg/simulator"
)

func doElevatorArrived(t *remoteController, msg *pb.SimulationEvent_ElevatorArrived) error {
	elevator := simulator.ElevatorID(msg.Arriving.ElevatorIndex)
	//dispatch to client
	t.controller.controller.CompletedMove(elevator)
	return nil
}
