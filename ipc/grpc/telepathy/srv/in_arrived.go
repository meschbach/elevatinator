package srv

import (
	"github.com/meschbach/elevatinator/ipc/grpc/telepathy/pb"
	"github.com/meschbach/elevatinator/simulator"
)

func doElevatorArrived(t *remoteController, msg *pb.SimulationEvent_ElevatorArrived) error {
	elevator := simulator.ElevatorID(msg.Arriving.ElevatorIndex)
	//dispatch to client
	t.controller.controller.CompletedMove(elevator)
	return nil
}
