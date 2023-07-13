package simulator

import "fmt"

type EventType int
type EntityID int
type FloorID int
type ElevatorID int

const (
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
	EventType EventType
	Timestamp Tick

	Entity   EntityID
	Elevator ElevatorID
	Floor    FloorID
	Points   int
}

type ControllerListener interface {
	OnControllerEvent(event Event)
}

func (e Event) ToString() string {
	switch e.EventType {
	case TickStart:
		return fmt.Sprintf("Event{TickStart,%d}", e.Timestamp)
	case TickDone:
		return fmt.Sprintf("Event{TickDone,%d}", e.Timestamp)
	case InitStart:
		return fmt.Sprintf("Event{InitStart}")
	case InitDone:
		return fmt.Sprintf("Event{InitDone}")
	case InformElevator:
		return fmt.Sprintf("Event{InformElevator, %d}", e.Elevator)
	case InformFloor:
		return fmt.Sprintf("Event{InformFloor, %d}", e.Floor)
	case ElevatorCalled:
		return fmt.Sprintf("Event{ElevatorCall, %d}", e.Floor)
	case ElevatorArrived:
		return fmt.Sprintf("Event{ElevatorArrived, %d @ %d}", e.Elevator, e.Floor)
	case ElevatorFloorRequest:
		return fmt.Sprintf("Event{ElevatorFloorRequest, %d @ %d}", e.Elevator, e.Floor)
	case ActorFinished:
		return fmt.Sprintf("Event{ActorFinished, point: %d}", e.Points)
	default:
		return fmt.Sprintf("Unkonwn event type %d: %#v", e.EventType, e)
	}
}

func OnTickStart(tick Tick) Event {
	return Event{
		EventType: TickStart,
		Timestamp: tick,
	}
}

func OnTickDone(tick Tick) Event {
	return Event{
		EventType: TickDone,
		Timestamp: tick,
	}
}

func OnInitStart() Event {
	return Event{
		EventType: InitStart,
		Timestamp: -1,
	}
}

func OnInitDone() Event {
	return Event{
		EventType: InitDone,
		Timestamp: -1,
	}
}

func OnInformElevator(id ElevatorID) Event {
	return Event{
		EventType: InformElevator,
		Timestamp: -1,
		Elevator:  id,
	}
}

func OnInformFloor(id FloorID) Event {
	return Event{
		EventType: InformFloor,
		Timestamp: -1,
		Floor:     id,
	}
}

func OnElevatorArrived(tick Tick, elevator ElevatorID, floor FloorID) Event {
	return Event{
		EventType: ElevatorArrived,
		Timestamp: tick,
		Elevator:  elevator,
		Floor:     floor,
	}
}

func OnElevatorCalled(tick Tick, floor FloorID) Event {
	return Event{
		EventType: ElevatorCalled,
		Timestamp: tick,
		Floor:     floor,
	}
}

func OnElevatorFloorRequest(tick Tick, elevator ElevatorID, floor FloorID) Event {
	return Event{
		EventType: ElevatorFloorRequest,
		Timestamp: tick,
		Elevator:  elevator,
		Floor:     floor,
	}
}

func OnActorFinished(tick Tick, points int) Event {
	return Event{
		EventType: ActorFinished,
		Timestamp: tick,
		Points:    points,
	}
}
