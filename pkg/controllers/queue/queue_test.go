package queue

import (
	"github.com/meschbach/elevatinator/pkg/scenarios"
	"testing"
)

func TestSingleUp(t *testing.T) {
	scenarios.TestScenario(t, NewController, scenarios.SinglePersonUp)
}

func TestSingleDown(t *testing.T) {
	scenarios.TestScenario(t, NewController, scenarios.SinglePersonDown)
}

func TestMultipleUpAndBack(t *testing.T) {
	scenarios.TestScenario(t, NewController, scenarios.MultipleUpAndBack)
}
