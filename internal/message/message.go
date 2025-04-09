package message

import (
	"encoding/json"
)

type Type string

const (
	TypeAuth      Type = "auth"
	TypeRaw            = "raw"
	TypeResponse       = "response"
	TypeTestLabel      = "test_label"
)

type Message struct {
	RequestID string          `json:"request_id,omitempty"`
	Type      Type            `json:"type"`
	Data      json.RawMessage `json:"data"`
}

type TestLabelData struct {
	PrinterID string `json:"printer_id"`
	Name      string `json:"name"`
	Street    string `json:"street"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Barcode   string `json:"barcode"`
}

type RawData struct {
	PrinterID string `json:"printer_id"`
	Data      string `json:"data"`
}

type ResponseData struct {
	RequestID string `json:"request_id"`
	Type      Type   `json:"type"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}

type AuthData struct {
	Name     string    `json:"name"`
	Key      string    `json:"key"`
	Printers []Printer `json:"printers"`
}

type Printer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type MessageOptionFunc func(*Message)

func New(typ Type, data interface{}, opts ...MessageOptionFunc) (*Message, error) {
	dataBytes, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	msg := &Message{
		Type: typ,
		Data: dataBytes,
	}

	for _, opt := range opts {
		opt(msg)
	}

	return msg, nil
}

func WithRequestID(v string) MessageOptionFunc {
	return func(m *Message) { m.RequestID = v }
}

func Decode(raw []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(raw, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func Encode(msg *Message) ([]byte, error) {
	d, err := json.Marshal(msg.Data)
	if err != nil {
		return nil, err
	}
	msg.Data = d
	return json.Marshal(&msg)
}
