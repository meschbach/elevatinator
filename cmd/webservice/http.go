package main

import (
	"context"
	"encoding/json"
	"net/http"
)

type httpReply interface {
	Write(w http.ResponseWriter) error
}

type JSONReply[T any] struct {
	Status int
	Body   T
}

func OkJSON[T any](value T) *JSONReply[T] {
	return &JSONReply[T]{
		Status: http.StatusOK,
		Body:   value,
	}
}

func (o *JSONReply[T]) Write(w http.ResponseWriter) error {
	header := w.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Content-Encoding", "UTF-8")
	w.WriteHeader(o.Status)
	return json.NewEncoder(w).Encode(o.Body)
}

func smartRoute(handler func(ctx context.Context, r *http.Request) (httpReply, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reply, err := handler(r.Context(), r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := reply.Write(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type ClientErrorReply struct {
	Problem error
}

func ClientError(problem error) *ClientErrorReply {
	return &ClientErrorReply{
		Problem: problem,
	}
}

func (c *ClientErrorReply) Write(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte(c.Problem.Error()))
	return err
}

func unprocessableEntity(message string) httpReply {
	return &staticReply{
		Status: http.StatusUnprocessableEntity,
		Body:   message,
	}
}

type staticReply struct {
	Status int
	Body   string
}

func (s *staticReply) Write(w http.ResponseWriter) error {
	w.WriteHeader(s.Status)
	_, err := w.Write([]byte(s.Body))
	return err
}

func notFound(message string) httpReply {
	return &staticReply{
		Status: http.StatusNotFound,
		Body:   message,
	}
}

type acceptedReply struct{}

var accepted = &acceptedReply{}

func (a *acceptedReply) Write(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusAccepted)
	return nil
}
