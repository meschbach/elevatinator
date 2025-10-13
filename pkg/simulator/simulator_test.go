package simulator

import (
	"fmt"
	"testing"
)

func TestSimulatorRunsForTicks(t *testing.T) {
	s := NewSimulation()
	s.Initialize(1, 2)
	for i := 0; i < 10; i++ {
		s.Tick()
	}
	if s.CurrentTick() != 10 {
		t.Errorf("Expect %d tickets, got %d", 10, s.CurrentTick())
	}
}

func TestActorAchievesGoal(t *testing.T) {
	targetFloor := 1
	maxTicks := Tick(10)
	s := NewSimulation()
	s.Initialize(1, 2)
	s.AttachActor(NewActor(targetFloor, 0, 0))
	s.AttachControllerFunc(NewMoveController)
	for i := Tick(0); i < maxTicks; i++ {
		s.Tick()
	}
	if s.actors[0].completedGoalTick == -1 {
		fmt.Printf("State: %#v", s.elevators[0])
		t.Errorf("Expected actor to achieve goal within %d ticks, took failed", maxTicks)
	}
}

func TestEventFeed(t *testing.T) {
	t.Logf("Starting event feed")
	capture := NewEventLog()
	maxTicks := Tick(10)
	s := NewSimulation()
	s.AttachActor(NewActor(1, 0, 0))
	s.AttachControllerListener(capture)

	t.Logf("initializing")
	s.Initialize(1, 2)
	s.AttachControllerFunc(NewMoveController)
	for i := Tick(0); i < maxTicks; i++ {
		s.Tick()
	}

	t.Log("Captured events")
	for _, e := range capture.Events {
		t.Log(e.ToString())
	}
}

func TestGameCompletion(t *testing.T) {
	capture := NewEventLog()
	maxTicks := Tick(10)
	s := NewSimulation()
	s.AttachActor(NewActor(1, 0, 0))
	s.AttachControllerListener(capture)

	t.Logf("initializing")
	s.Initialize(1, 2)
	s.AttachControllerFunc(NewMoveController)
	endTick := s.TickUpTo(maxTicks)

	if endTick >= maxTicks {
		t.Errorf("Exceeded tick count @ %d", s.tick)
	}
}

func TestActorsNotCompletedAtStart(t *testing.T) {
	s := NewSimulation()
	s.AttachActor(NewActor(1, 0, 0))

	t.Logf("initializing")
	if s.ActorsCompletedObjectives() {
		t.Errorf("Game completed before initialized")
	}
}

func TestElevatorsMoveDown(t *testing.T) {
	maxTicks := Tick(10)
	s := NewSimulation()
	s.AttachActor(NewActor(0, 1, 0))

	s.Initialize(1, 2)
	s.AttachControllerFunc(NewMoveController)
	endTick := s.TickUpTo(maxTicks)

	if endTick >= maxTicks {
		t.Errorf("Exceeded tick count @ %d", s.tick)
	}
}
