package main

import (
	"errors"
	"sync"

	"github.com/meschbach/elevatinator/pkg/simulator"
	"github.com/meschbach/go-junk-bucket/pkg/fx"
)

type gameSession struct {
	state      sync.RWMutex
	simulation *simulator.Simulation
	isDone     bool
	eventLog   *gameSessionLog
}

func (s *service) newGameSession(scenarioName string, aiName string) (*gameSession, error) {
	matchedAIUnits := fx.Filter(s.aiUnits, func(aiUnits aiUnits) bool {
		return aiUnits.Name == aiName
	})
	matchedScenario := fx.Filter(s.availableScenarios, func(scenario scenario) bool {
		return scenario.Name == scenarioName
	})
	if len(matchedAIUnits) != 1 {
		return nil, errors.New("no matching ai unit")
	}
	if len(matchedScenario) != 1 {
		return nil, errors.New("no matching scenario")
	}

	log := &gameSessionLog{}

	sim := simulator.NewSimulation()
	sim.AttachControllerListener(log)
	matchedScenario[0].setup(sim)
	sim.AttachControllerFunc(matchedAIUnits[0].Controller)

	return &gameSession{
		state:      sync.RWMutex{},
		simulation: sim,
		isDone:     false,
		eventLog:   log,
	}, nil
}

func (g *gameSession) tick() (done bool, problem error) {
	g.state.Lock()
	defer g.state.Unlock()

	if g.isDone {
		return true, nil
	}

	g.isDone = !g.simulation.Tick()
	return g.isDone, nil
}

type gameSessionLogMarker int

type gameSessionLog struct {
	state  sync.RWMutex
	events []simulator.Event
}

func (g *gameSessionLog) OnControllerEvent(event simulator.Event) {
	g.state.Lock()
	defer g.state.Unlock()

	g.events = append(g.events, event)
}

func (g *gameSessionLog) Marker() gameSessionLogMarker {
	g.state.Lock()
	defer g.state.Unlock()

	return gameSessionLogMarker(len(g.events))
}

func (g *gameSessionLog) EventsSince(marker gameSessionLogMarker) []simulator.Event {
	g.state.RLock()
	defer g.state.RUnlock()

	return g.events[marker:]
}

func (g *gameSessionLog) Events() []simulator.Event {
	g.state.Lock()
	defer g.state.Unlock()

	output := make([]simulator.Event, len(g.events))
	copy(output, g.events)
	return output
}
