package srv

import (
	"github.com/meschbach/elevatinator/pkg/ipc/grpc/telepathy/pb"
	"github.com/meschbach/elevatinator/pkg/simulator"
)

func doFloorCall(t *remoteController, msg *pb.SimulationEvent_ElevatorCalled) error {
	id := simulator.FloorID(msg.CalledAt.FloorIndex)
	//dispatch to client
	t.controller.controller.Called(id)
	return nil
}
