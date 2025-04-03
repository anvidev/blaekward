package event

import "encoding/json"

type Event struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

type PrintPayload struct {
	Addr string `json:"addr"`
	Data string `json:"data"`
}

type EventHandler interface {
	Handle(Event) error
}

// event handler interface for both server and client
// implementation can have a map of map[Event]func(Event) error
// for easy event handling

type EEvent interface {
	Encode() ([]byte, error)
	Decode(data []byte) error
	EventType() string
}

type HandleFunc func(EEvent) error

type PrintEvent struct {
	PrinterID string `json:"printer_id"`
	Data      string `json:"data"`
}

func (e *PrintEvent) Encode() ([]byte, error) {
	return json.Marshal(e)
}
func (e *PrintEvent) Decode(data []byte) error {
	return json.Unmarshal(data, e)
}
func (e *PrintEvent) EventType() string { return "Print" }
