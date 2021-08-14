package srv

import (
	"fmt"
	"github.com/meschbach/elevatinator/simulator"
)

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