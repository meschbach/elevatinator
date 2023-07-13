package simulator

type EventLog struct {
	Events []Event
}

func (e *EventLog) OnControllerEvent(event Event) {
	e.Events = append(e.Events, event)
}

func NewEventLog() *EventLog {
	return &EventLog{Events: make([]Event, 0)}
}
