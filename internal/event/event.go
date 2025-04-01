package event

import "encoding/json"

type Event struct {
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload"`
}

type TestPayload struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PrintPayload struct {
	IP   string `json:"ip"`
	Data string `json:"data"`
}
