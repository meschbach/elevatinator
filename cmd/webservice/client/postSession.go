package main

import "context"

type PostSessionRequestBody struct {
	Scenario   *string `json:"scenario,omitempty"`
	Controller *string `json:"controller,omitempty"`
}

type PostSessionResponseBody struct {
	SessionID string `json:"sessionID"`
}

func (c *webClient) PostSession(ctx context.Context, body PostSessionRequestBody) (*PostSessionResponseBody, error) {
	req, err := c.NewRequest(ctx, "POST", "/session", body)
	if err != nil {
		return nil, err
	}

	out := &PostSessionResponseBody{}
	if err := c.Do(req, &out); err != nil {
		return nil, err
	}
	return out, nil
}
