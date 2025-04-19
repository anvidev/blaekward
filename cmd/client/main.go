package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/anvidev/blaekward/internal/event"
	ws "github.com/gorilla/websocket"
)

type Gateway struct {
	logger *slog.Logger

	mu    sync.Mutex
	conn  *ws.Conn
	setup setup
}

type setup struct {
	Name     string                   `json:"name"`
	Key      string                   `json:"key"`
	Printers map[string]event.Printer `json:"printers"`
}

func main() {
	gw := &Gateway{logger: slog.Default()}

	if err := gw.connect(); err != nil {
		gw.logger.Error("connecting to server failed", "error", err)
		return
	}
	defer gw.conn.Close()

	if err := gw.authenticate(); err != nil {
		gw.logger.Error("authentication failed on server", "error", err)
		return
	}

	gw.read()
}

func (gw *Gateway) connect() error {
	if err := gw.loadIniFile(); err != nil {
		return err
	}

	url := url.URL{
		Scheme: "ws",
		Host:   "localhost:9090",
		Path:   "/v1/ws",
	}

	conn, _, err := ws.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return err
	}

	gw.logger.Info("connected to server")
	gw.conn = conn

	return nil
}

func (gw *Gateway) authenticate() error {
	gw.conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	defer gw.conn.SetReadDeadline(time.Time{})

	authData := event.AuthData{
		Name:     gw.setup.Name,
		Key:      gw.setup.Key,
		Printers: []event.Printer{},
	}

	for id, printer := range gw.setup.Printers {
		printer.ID = id
		authData.Printers = append(authData.Printers, printer)
	}

	data, err := json.Marshal(authData)
	if err != nil {
		return err
	}

	evt := event.New(event.TypeAuth, data)

	if err := gw.conn.WriteJSON(evt); err != nil {
		return err
	}

	var resp event.Event
	if err := gw.conn.ReadJSON(&resp); err != nil {
		return err
	}

	switch resp.Type {
	case event.TypeError:
		var errorData event.ErrorData
		if err := event.Extract(resp.Data, &errorData); err != nil {
			return err
		}
		if errorData.Fatal {
			return fmt.Errorf(errorData.Message)
		}
	}

	return nil
}

func (gw *Gateway) read() error {
	for {
		var evt event.Event
		if err := gw.conn.ReadJSON(&evt); err != nil {
			// gw.logger.Error("read error", "error", err)
			return err
		}

		switch evt.Type {
		case event.TypeRaw:
			var rawData event.RawData
			if err := event.Extract(evt.Data, &rawData); err != nil {
				gw.logger.Warn("error unmarshaling raw event type data", "error", err)
				continue
			}
			gw.handleRaw(rawData, evt.RequestID)
		case event.TypeError:
			var errorData event.ErrorData
			if err := event.Extract(evt.Data, &errorData); err != nil {
				gw.logger.Warn("error unmarshaling error event  type data", "error", err)
				continue
			}
			gw.logger.Error("error received from server", "event", evt.Type, "request_id", evt.RequestID)
			if errorData.Fatal {
				return fmt.Errorf(errorData.Message)
			}
		default:
			gw.logger.Warn("unknown event received", "event", evt.Type)
		}
	}
}

func (gw *Gateway) handleRaw(raw event.RawData, reqID string) {
	gw.mu.Lock()
	defer gw.mu.Unlock()

	printer, found := gw.setup.Printers[raw.PrinterID]
	if !found {
		data, err := json.Marshal(event.ErrorData{
			Code:    event.CodeNotFound,
			Message: fmt.Sprintf("Printer with id %s was not found", raw.PrinterID),
		})
		if err != nil {
			gw.logger.Error("error marshaling error response data", "error", err)
			return
		}
		evt := event.New(event.TypeError, data, event.WithRequestID(reqID))
		gw.conn.WriteJSON(evt)
		return
	}

	addr := fmt.Sprintf("%s:%s", printer.Addr, "9100")
	if err := gw.sendJSONToPrinter(addr, raw); err != nil {
		data, err := json.Marshal(event.ErrorData{
			Code:    event.CodeNotReachable,
			Message: err.Error(),
		})
		if err != nil {
			gw.logger.Error("error marshaling error response data", "error", err)
			return
		}
		evt := event.New(event.TypeError, data, event.WithRequestID(reqID))
		gw.conn.WriteJSON(evt)
		return
	}

	data, err := json.Marshal(event.SuccessData{
		Message: "Event was successful",
	})
	if err != nil {
		gw.logger.Error("error marshaling success response data", "error", err)
		return
	}
	evt := event.New(event.TypeSuccess, data, event.WithRequestID(reqID))
	if err := gw.conn.WriteJSON(evt); err != nil {
		gw.logger.Error("error writing json success response data", "error", err)
		return
	}
}

func (gw *Gateway) sendJSONToPrinter(addr string, data any) error {
	tcp, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp4", nil, tcp)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
