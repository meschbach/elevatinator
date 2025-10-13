package main

import "context"

type GetControllersDescription struct {
	Name string
}

type GetControllersReply struct {
	Available []GetControllersDescription `json:"available"`
}

func (c *webClient) GetControllers(ctx context.Context) (*GetControllersReply, error) {
	req, err := c.NewRequest(ctx, "GET", "/controllers", nil)
	if err != nil {
		return nil, err
	}

	out := &GetControllersReply{}
	if err := c.Do(req, &out); err != nil {
		return nil, err
	}
	return out, nil
}
