package event

import (
	"encoding/json"
)

type EventType string

type ErrorCode string

const (
	CodeUnauthorized ErrorCode = "UNAUTHORIZED"
	CodeBadRequest             = "BAD_REQUEST"
	CodeNotFound               = "NOT_FOUND"
	CodeNotReachable           = "NOT_REACHABLE"
	CodeWriteFailed            = "WRITE_FAILED"

	TypeError   EventType = "error"
	TypeSuccess           = "success"
	TypeRaw               = "raw"
	TypeAuth              = "auth"
)

type EventOption func(*Event)

// Event is the data structure for communicating between server and client.
//
// Both server and client handles events based in the event type, and only after
// determining the type should the data (or payload) get processed, as this is
// just a slice of bytes.
//
// Events also carry the request id from the initial api request, which the client
// must include in the respoinding event to allow the server to distribute the event
// to the correct api resquest in order for the request to respond to the consumer.
type Event struct {
	Type EventType       `json:"type"`
	Data json.RawMessage `json:"data"`

	RequestID string `json:"request_id"`
}

type RawData struct {
	PrinterID string `json:"printer_id"`
	Data      string `json:"data"`
}

type SuccessData struct {
	Message string `json:"message"`
}

type ErrorData struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Fatal   bool      `json:"fatal"`
}

type AuthData struct {
	Name     string    `json:"name"`
	Key      string    `json:"key"`
	Printers []Printer `json:"printers"`
}

// Printer is the helper struct that both the server and client uses in their
// internal state.
type Printer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Addr string `json:"addr"`
}

// New initializes an Event ready to be sent over the connection.
func New(typ EventType, data json.RawMessage, opts ...EventOption) Event {
	evt := Event{Type: typ, Data: data}

	for _, opt := range opts {
		opt(&evt)
	}

	return evt
}

func WithRequestID(id string) EventOption {
	return func(e *Event) { e.RequestID = id }
}

// Extract unmarshals the incoming data into [T] which should be any of the data structs.
func Extract[T any](data []byte, v *T) error {
	return json.Unmarshal(data, v)
}
