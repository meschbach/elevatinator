package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/meschbach/elevatinator/pkg/simulator"
)

type DynamicScenarioWire struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Floors      int                        `json:"floors"`
	Elevators   int                        `json:"elevators"`
	Actors      []DynamicScenarioActorWire `json:"actors"`
}

func (d *DynamicScenarioWire) validate() string {
	if d.Name == "" {
		return "missing name"
	}
	if d.Floors < 1 {
		return "floors must be at least 1"
	}
	if d.Elevators < 1 {
		return "elevators must be at least 1"
	}
	return ""
}

type DynamicScenarioActorWire struct {
	Name          string            `json:"name"`
	StartingFloor simulator.FloorID `json:"starting-floor"`
	StartingTick  simulator.Tick    `json:"starting-tick"`
	GoalFloor     simulator.FloorID `json:"goal-floor"`
}

type PostScenarioRequestBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PostScenarioResponseBody struct {
	ID string `json:"id"`
}

func (s *service) postScenarioRoute(ctx context.Context, r *http.Request) (httpReply, error) {
	var requestBody PostScenarioRequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return ClientError(err), nil
	}
	if requestBody.Name == "" {
		return unprocessableEntity("missing name"), nil
	}

	//generate ID
	generatedUUID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	id := generatedUUID.String()

	//parse method body
	s.state.Lock()
	defer s.state.Unlock()
	s.dynamicScenarios[id] = &DynamicScenarioWire{
		ID:          id,
		Name:        requestBody.Name,
		Description: requestBody.Description,
		Floors:      0,
		Elevators:   0,
		Actors:      nil,
	}

	//
	return OkJSON(PostScenarioResponseBody{
		ID: id,
	}), nil
}

type GetScenarioResponseBody struct {
	Scenarios []GetScenarioResponseBodyScenario `json:"scenarios"`
}
type GetScenarioResponseBodyScenario struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *service) getScenarioRoute(ctx context.Context, r *http.Request) (httpReply, error) {
	//todo: could tighten up the locking here
	s.state.RLock()
	defer s.state.RUnlock()

	output := make([]GetScenarioResponseBodyScenario, len(s.dynamicScenarios))
	for _, scenario := range s.dynamicScenarios {
		output = append(output, GetScenarioResponseBodyScenario{
			ID:          scenario.ID,
			Name:        scenario.Name,
			Description: scenario.Description,
		})
	}
	return OkJSON(GetScenarioResponseBody{
		Scenarios: output,
	}), nil
}

func (s *service) putScenarioRoute(ctx context.Context, r *http.Request) (httpReply, error) {
	scenarioID := mux.Vars(r)["scenarioID"]
	//Parse body
	var requestBody DynamicScenarioWire
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return ClientError(err), nil
	}
	if validationErrors := requestBody.validate(); validationErrors != "" {
		return unprocessableEntity(validationErrors), nil
	}

	// lock the body
	s.state.Lock()
	defer s.state.Unlock()

	if _, ok := s.dynamicScenarios[requestBody.ID]; !ok {
		return notFound(fmt.Sprintf("scenario %q not found", requestBody.ID)), nil
	}

	s.dynamicScenarios[scenarioID] = &requestBody
	return accepted, nil
}

func (s *service) deleteScenarioRoute(ctx context.Context, r *http.Request) (httpReply, error) {
	s.state.Lock()
	defer s.state.Unlock()
	delete(s.dynamicScenarios, mux.Vars(r)["scenarioID"])
	return accepted, nil
}
