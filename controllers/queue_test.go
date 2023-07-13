package controllers

import (
	"github.com/meschbach/elevatinator/pkg/scenarios"
	"testing"
)

func TestSingleUp(t *testing.T) {
	scenarios.TestScenario(t, NewQueueController, scenarios.SinglePersonUp)
}

func TestSingleDown(t *testing.T) {
	scenarios.TestScenario(t, NewQueueController, scenarios.SinglePersonDown)
}

func TestMultipleUpAndBack(t *testing.T) {
	scenarios.TestScenario(t, NewQueueController, scenarios.MultipleUpAndBack)
}
