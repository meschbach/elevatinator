package main

import (
	"context"
	"fmt"
	"net/http"
)

type PostSessionTickResponseBody struct {
	Completed bool `json:"completed"`
}

func (c *webClient) PostSessionTick(ctx context.Context, sessionID string) (*PostSessionTickResponseBody, error) {
	url := fmt.Sprintf("/session/%s/tick", sessionID)
	req, err := c.NewRequest(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	out := &PostSessionTickResponseBody{}
	if err := c.Do(req, out); err != nil {
		return nil, err
	}
	return out, nil
}
