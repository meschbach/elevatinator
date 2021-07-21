package main

import (
	"github.com/meschbach/elevatinator/scenarios"
	"github.com/meschbach/elevatinator/simulator"
)

func main() {
	scenarios.RunScenario(simulator.NewMoveController, scenarios.MultipleUpAndBack)
}
