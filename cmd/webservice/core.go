package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/meschbach/elevatinator/pkg/simulator"
)

type scenario struct {
	Name        string
	Description string
	setup       func(simulation *simulator.Simulation) simulator.Tick
}

type aiUnits struct {
	Name       string
	Controller simulator.ControllerFunc
}

type service struct {
	builtinScenarios []scenario
	aiUnits          []aiUnits
	dynamicScenarios map[string]*DynamicScenarioWire

	state        *sync.RWMutex
	gameSessions map[string]*gameSession
}

type GetScenariosDescription struct {
	Name        string
	Description string
}

type GetScenariosReply struct {
	Available []GetScenariosDescription `json:"available"`
}

func (s *service) getScenariosRoute(ctx context.Context, r *http.Request) (httpReply, error) {
	output := GetScenariosReply{}
	for _, s := range s.builtinScenarios {
		output.Available = append(output.Available, GetScenariosDescription{
			Name:        s.Name,
			Description: s.Description,
		})
	}

	for _, s := range s.dynamicScenarios {
		output.Available = append(output.Available, GetScenariosDescription{
			Name:        s.Name,
			Description: s.Description,
		})
	}

	return OkJSON(output), nil
}

type GetControllersDescription struct {
	Name string
}

type GetControllersReply struct {
	Available []GetControllersDescription `json:"available"`
}

func (s *service) getControllersRoute(ctx context.Context, r *http.Request) (httpReply, error) {
	output := GetControllersReply{}
	for _, s := range s.aiUnits {
		output.Available = append(output.Available, GetControllersDescription{
			Name: s.Name,
		})
	}

	return OkJSON(output), nil
}

type PostSessionRequestBody struct {
	Scenario   *string `json:"scenario,omitempty"`
	Controller *string `json:"controller,omitempty"`
}

type PostSessionResponseBody struct {
	SessionID string `json:"sessionID"`
}

func (s *service) postSessionRoute(ctx context.Context, r *http.Request) (httpReply, error) {
	//
	requestBody := PostSessionRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return ClientError(err), nil
	}
	if requestBody.Scenario == nil {
		return unprocessableEntity("missing scenario"), nil
	}
	if requestBody.Controller == nil {
		return unprocessableEntity("missing controller"), nil
	}

	//build the session
	session, err := s.newGameSession(*requestBody.Scenario, *requestBody.Controller)
	if err != nil {
		return nil, err
	}

	// generate ID
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	idString := id.String()

	// attach to the service
	func() {
		s.state.Lock()
		defer s.state.Unlock()
		s.gameSessions[idString] = session
	}()

	// return the session ID
	return OkJSON(PostSessionResponseBody{
		SessionID: idString,
	}), nil
}

type PostSessionTickResponseBody struct {
	Completed bool `json:"completed"`
}

func (s *service) postSessionTickRoute(ctx context.Context, r *http.Request) (httpReply, error) {
	// Grab the session
	pathVariables := mux.Vars(r)
	id := pathVariables["sessionID"]

	session := func() *gameSession {
		s.state.RLock()
		defer s.state.RUnlock()

		session := s.gameSessions[id]
		return session
	}()

	if session == nil {
		return notFound(fmt.Sprintf("session %q not found", id)), nil
	}

	// have the session tick forward
	if completed, err := session.tick(); err != nil {
		return nil, err
	} else {
		return OkJSON(PostSessionTickResponseBody{Completed: completed}), nil
	}
}

type GetSessionEventsReply struct {
	Events []GetSessionEventsReplyEvents `json:"events"`
}

type GetSessionEventsReplyEvents struct {
	EventType string          `json:"eventType"`
	Timestamp *simulator.Tick `json:"timestamp,omitempty"`

	Entity   *simulator.EntityID   `json:"entity,omitempty"`
	Elevator *simulator.ElevatorID `json:"elevator,omitempty"`
	Floor    *simulator.FloorID    `json:"floor,omitempty"`
	Points   *int                  `json:"points,omitempty"`
}

func (s *service) getSessionEvents(ctx context.Context, r *http.Request) (httpReply, error) {
	// Grab the session
	pathVariables := mux.Vars(r)
	id := pathVariables["sessionID"]

	session := func() *gameSession {
		s.state.RLock()
		defer s.state.RUnlock()

		session := s.gameSessions[id]
		return session
	}()

	if session == nil {
		return notFound(fmt.Sprintf("session %q not found", id)), nil
	}

	// grab the session logs
	events := session.eventLog.Events()
	translated := make([]GetSessionEventsReplyEvents, len(events))
	for index, event := range events {
		switch event.EventType {
		case simulator.TickStart:
			translated[index].EventType = "TickStart"
			translated[index].Timestamp = &event.Timestamp
		case simulator.TickDone:
			translated[index].EventType = "TickDone"
			translated[index].Timestamp = &event.Timestamp
		case simulator.InitStart:
			translated[index].EventType = "InitStart"
		case simulator.InitDone:
			translated[index].EventType = "InitDone"
		case simulator.InformElevator:
			translated[index].EventType = "InformElevator"
			translated[index].Elevator = &event.Elevator
		case simulator.InformFloor:
			translated[index].EventType = "InformFloor"
			translated[index].Floor = &event.Floor
		case simulator.ElevatorCalled:
			translated[index].EventType = "ElevatorCalled"
			translated[index].Floor = &event.Floor
		case simulator.ElevatorArrived:
			translated[index].EventType = "ElevatorArrived"
			translated[index].Floor = &event.Floor
			translated[index].Elevator = &event.Elevator
		case simulator.ElevatorFloorRequest:
			translated[index].EventType = "ElevatorFloorRequest"
			translated[index].Floor = &event.Floor
			translated[index].Elevator = &event.Elevator
		case simulator.ActorFinished:
			translated[index].EventType = "ActorFinished"
			translated[index].Points = &event.Points
		case simulator.ElevatorAtFloor:
			translated[index].EventType = "ElevatorAtFloor"
			translated[index].Floor = &event.Floor
			translated[index].Elevator = &event.Elevator
		default:
			translated[index].EventType = fmt.Sprintf("%s", event.ToString())
		}
	}

	return OkJSON(GetSessionEventsReply{
		Events: translated,
	}), nil
}
