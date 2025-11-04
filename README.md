# Elevatinator

Build an Elevator Controller and compete against your friends!  See who can build the most optimal controller.

## Getting Started

Checkout [main.go](main.go)!  This will build a scenario and execute the specified controller in source.  If all actors
get to their floor it will print out how many ticks it took to get there!  See if your friends algorithms are faster.

Build your own controller and replace the following call with a function to instantiate it:
```go
simulation.AttachControllerFunc(simulator.NewMoveController)
```

### New to Go?

In the future I would like to add other language via an RPC/IPC mechanism.  For now one can create a file in the same
directory as `main.go`, called `controller.go`.  Use the following template to get you started:

```go
package main

import "github.com/meschbach/elevatinator/pkg/simulator"

type MyStrategy struct {
	// any data or state should go here
}

func NewStrategy(elevators simulator.ControlledElevators) simulator.Controller {
	return &MyStrategy{}
}

func (m *MyStrategy) Init(elevators []simulator.ElevatorID)                                  {}
func (m *MyStrategy) Called(floor simulator.FloorID)                                         {}
func (m *MyStrategy) FloorSelected(elevatorID simulator.ElevatorID, floor simulator.FloorID) {}
func (m *MyStrategy) CompletedMove(elevatorID simulator.ElevatorID)                          {}
```

This allows you to plugin to the simulation.  Additionally, you'll need to modify `main.go` from
`scenarios.RunScenario(simulator.NewMoveController, scenarios.MultipleUpAndBack)` *to* `scenarios.RunScenario(NewStrategy, scenarios.MultipleUpAndBack)`

Check out [simulator/movecontroller.go](pkg/simulator/movecontroller.go) for  examples on how to move elevators!

#### Building & Running

Once you are ready to test your new strategy can run it via the following:

```bash
go build . && ./elevatinator
```
