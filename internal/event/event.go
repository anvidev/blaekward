package event

import "encoding/json"

type Event struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

type PrintPayload struct {
	PrinterID string `json:"printer_id"`
	Data      string `json:"data"`
}

type IEvent interface {
	Encode() ([]byte, error)
	Decode(data []byte) error
	Event() string
}

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
func (e *PrintEvent) Event() string { return "Print" }

type Printer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type AuthenticationEvent struct {
	Name     string    `json:"name"`
	CID      string    `json:"cid"`
	Key      string    `json:"key"`
	Printers []Printer `json:"printers"`
}

func (e *AuthenticationEvent) Encode() ([]byte, error) {
	return json.Marshal(e)
}
func (e *AuthenticationEvent) Decode(data []byte) error {
	return json.Unmarshal(data, e)
}
func (e *AuthenticationEvent) EventType() string { return "Authentication" }
