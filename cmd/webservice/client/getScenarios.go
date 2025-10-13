package main

import "context"

type GetScenariosDescription struct {
	Name        string
	Description string
}

type GetScenariosReply struct {
	Available []GetScenariosDescription `json:"available"`
}

func (c *webClient) GetScenarios(ctx context.Context) (*GetScenariosReply, error) {
	req, err := c.NewRequest(ctx, "GET", "/scenarios", nil)
	if err != nil {
		return nil, err
	}

	out := &GetScenariosReply{}
	if err := c.Do(req, &out); err != nil {
		return nil, err
	}
	return out, nil
}
