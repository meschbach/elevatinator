package srv

import (
	"fmt"
	simulator2 "github.com/meschbach/elevatinator/pkg/simulator"
)

type pendingMove struct {
	which simulator2.ElevatorID
	to    simulator2.FloorID
}

type controllerInstance struct {
	controller   simulator2.Controller
	pending      []*pendingMove
	maxElevators uint32
}

func (c *controllerInstance) MoveTo(elevator simulator2.ElevatorID, floor simulator2.FloorID) {
	if elevator < 0 || uint32(elevator) > c.maxElevators {
		//TODO: Report problem
		panic(fmt.Sprintf("no such elevator %d", elevator))
	}
	fmt.Printf("Queuing move of %d to %d\n", elevator, floor)
	c.pending = append(c.pending, &pendingMove{
		which: elevator,
		to:    floor,
	})
}

func (c *controllerInstance) resetPending() {
	c.pending = make([]*pendingMove, 0)
}
