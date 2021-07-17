package simulator

import "fmt"

type EventType int
type EntityID int
type FloorID int
type ElevatorID int

const (
	//No operands
	TickStart = iota
	TickDone

	InitStart
	InformElevator
	InformFloor
	InitDone

	ElevatorCalled
	ElevatorArrived
	ElevatorFloorRequest

	ActorFinished
)

type Event struct {
	eventType EventType
	timestamp Tick

	entity   EntityID
	elevator ElevatorID
	floor    FloorID
	points   int
}

type SimulatorControllerListener interface {
	OnControllerEvent(event Event)
}

func (e Event) ToString() string {
	switch e.eventType {
	case TickStart:
		return fmt.Sprintf("Event{TickStart,%d}", e.timestamp)
	case TickDone:
		return fmt.Sprintf("Event{TickDone,%d}", e.timestamp)
	case InitStart:
		return fmt.Sprintf("Event{InitStart}")
	case InitDone:
		return fmt.Sprintf("Event{InitDone}")
	case InformElevator:
		return fmt.Sprintf("Event{InformElevator, %d}", e.elevator)
	case InformFloor:
		return fmt.Sprintf("Event{InformFloor, %d}", e.floor)
	case ElevatorCalled:
		return fmt.Sprintf("Event{ElevatorCall, %d}", e.floor)
	case ElevatorArrived:
		return fmt.Sprintf("Event{ElevatorArrived, %d @ %d}", e.elevator, e.floor)
	case ElevatorFloorRequest:
		return fmt.Sprintf("Event{ElevatorFloorRequest, %d @ %d}", e.elevator, e.floor)
	case ActorFinished:
		return fmt.Sprintf("Event{ActorFinished, point: %d}", e.points)
	default:
		return fmt.Sprintf("Unkonwn event type %d: %#v", e.eventType, e)
	}
}

func OnTickStart(tick Tick) Event {
	return Event{
		eventType: TickStart,
		timestamp: tick,
	}
}

func OnTickDone(tick Tick) Event {
	return Event{
		eventType: TickDone,
		timestamp: tick,
	}
}

func OnInitStart() Event {
	return Event{
		eventType: InitStart,
		timestamp: -1,
	}
}

func OnInitDone() Event {
	return Event{
		eventType: InitDone,
		timestamp: -1,
	}
}

func OnInformElevator(id ElevatorID) Event {
	return Event{
		eventType: InformElevator,
		timestamp: -1,
		elevator:  id,
	}
}

func OnInformFloor(id FloorID) Event {
	return Event{
		eventType: InformFloor,
		timestamp: -1,
		floor:     id,
	}
}

func OnElevatorArrived(tick Tick, elevator ElevatorID, floor FloorID) Event {
	return Event{
		eventType: ElevatorArrived,
		timestamp: tick,
		elevator:  elevator,
		floor:     floor,
	}
}

func OnElevatorCalled(tick Tick, floor FloorID) Event {
	return Event{
		eventType: ElevatorCalled,
		timestamp: tick,
		floor:     floor,
	}
}

func OnElevatorFloorRequest(tick Tick, elevator ElevatorID, floor FloorID) Event {
	return Event{
		eventType: ElevatorFloorRequest,
		timestamp: tick,
		elevator:  elevator,
		floor:     floor,
	}
}

func OnActorFinished(tick Tick, points int) Event {
	return Event{
		eventType: ActorFinished,
		timestamp: tick,
		points: points,
	}
}
