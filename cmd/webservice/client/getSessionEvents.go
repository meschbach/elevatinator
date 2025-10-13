package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/meschbach/elevatinator/pkg/simulator"
)

type GetSessionEventsReply struct {
	Events []GetSessionEventsReplyEvents `json:"events"`
}

type GetSessionEventsReplyEvents struct {
	EventType string `json:"eventType"`
	Timestamp *int64 `json:"timestamp,omitempty"`

	Entity   *simulator.EntityID   `json:"entity,omitempty"`
	Elevator *simulator.ElevatorID `json:"elevator,omitempty"`
	Floor    *simulator.FloorID    `json:"floor,omitempty"`
	Points   *int                  `json:"points,omitempty"`
}

func (c *webClient) GetSessionEvents(ctx context.Context, sessionID string) (*GetSessionEventsReply, error) {
	url := fmt.Sprintf("/session/%s/events", sessionID)
	req, err := c.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	out := &GetSessionEventsReply{}
	if err := c.Do(req, out); err != nil {
		fmt.Printf("Error invoking %s: %s\n", url, err.Error())
		return nil, err
	}
	return out, nil
}
