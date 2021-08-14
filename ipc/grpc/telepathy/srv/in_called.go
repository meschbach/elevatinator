package srv

import (
	"github.com/meschbach/elevatinator/ipc/grpc/telepathy/pb"
	"github.com/meschbach/elevatinator/simulator"
)

func doFloorCall(t *remoteController, msg *pb.SimulationEvent_ElevatorCalled) error {
	id := simulator.FloorID(msg.CalledAt.FloorIndex)
	//dispatch to client
	t.controller.controller.Called(id)
	return nil
}
