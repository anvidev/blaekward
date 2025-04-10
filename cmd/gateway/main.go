package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/anvidev/blaekward/internal/conn"
	"github.com/anvidev/blaekward/internal/label"
	"github.com/anvidev/blaekward/internal/message"
)

type Gateway struct {
	logger *slog.Logger
	conn   *conn.WrappedConn
	State  state
}

type state struct {
	Key      string                     `json:"key"`
	Name     string                     `json:"name"`
	Printers map[string]message.Printer `json:"printers"`
}

func main() {
	gw := &Gateway{logger: slog.Default()}

	blaekwardUrl := "blaekward.anvi.dev:9091"

	if err := gw.connect(blaekwardUrl); err != nil {
		gw.logger.Error("connection failed", "error", err)
		return
	}

	gw.listen()
}

func (gw *Gateway) listen() {
	for {
		data, err := gw.conn.ReadMessage()
		if err != nil {
			gw.logger.Error("read error", "error", err)
			return
		}

		msg, err := message.Decode(data)
		if err != nil {
			gw.logger.Error("decode error", "error", err)
			continue
		}

		gw.logger.Info("received message from server", "type", msg.Type)

		switch msg.Type {
		case message.TypeRaw:
			gw.handleRaw(msg)
		case message.TypeTestLabel:
			gw.handleTestLabel(msg)
		default:
			gw.logger.Warn("unhandled message type", "type", msg.Type)
			continue
		}
	}
}

func (gw *Gateway) handleTestLabel(msg *message.Message) {
	var data message.TestLabelData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		gw.logger.Error("failed to decode message", "type", msg.Type, "error", err)
		return
	}

	printer, found := gw.State.Printers[data.PrinterID]
	if !found {
		res := &message.ResponseData{
			RequestID: msg.RequestID,
			Type:      msg.Type,
			Success:   false,
			Message:   "printer not found",
		}

		respMsg, err := message.New(message.TypeResponse, res)
		if err != nil {
			gw.logger.Error("failed to encode response message", "error", err)
			return
		}

		if err := gw.sendResponse(respMsg); err != nil {
			gw.logger.Error("failed to send response", "error", err)
			return
		}
	}

	addr := fmt.Sprintf("%s:%s", printer.Addr, "9100")
	printerConn, err := dialPrinter(addr)
	if err != nil {
		gw.logger.Error("failed to dial printer", "error", err)

		res := &message.ResponseData{
			RequestID: msg.RequestID,
			Type:      msg.Type,
			Success:   false,
			Message:   fmt.Sprintf("failed to dial printer %s on %s", printer.Name, addr),
		}

		respMsg, err := message.New(message.TypeResponse, res)
		if err != nil {
			gw.logger.Error("failed to encode response message", "error", err)
			return
		}

		if err := gw.sendResponse(respMsg); err != nil {
			gw.logger.Error("failed to send response", "error", err)
			return
		}
	}

	labelBytes := label.NewTest(
		data.Name,
		data.Street,
		data.City,
		data.Country,
		data.Barcode,
	)

	if _, err := printerConn.Write(labelBytes); err != nil {
		res := &message.ResponseData{
			RequestID: msg.RequestID,
			Type:      msg.Type,
			Success:   false,
			Message:   fmt.Sprintf("failed to write to printer %s on %s", printer.Name, addr),
		}

		respMsg, err := message.New(message.TypeResponse, res)
		if err != nil {
			gw.logger.Error("failed to encode response message", "error", err)
			return
		}

		if err := gw.sendResponse(respMsg); err != nil {
			gw.logger.Error("failed to send response", "error", err)
			return
		}
	}

	response := &message.ResponseData{
		RequestID: msg.RequestID,
		Type:      msg.Type,
		Message:   "message handled successfully",
		Success:   true,
	}

	respMsg, err := message.New(message.TypeResponse, response)
	if err != nil {
		gw.logger.Error("failed to encode response message", "error", err)
	}

	if err := gw.sendResponse(respMsg); err != nil {
		gw.logger.Error("failed to send response", "error", err)
	}
}

func (gw *Gateway) handleRaw(msg *message.Message) {
	var msgData message.RawData
	if err := json.Unmarshal(msg.Data, &msgData); err != nil {
		gw.logger.Error("failed to decode message", "type", msg.Type, "error", err)
		return
	}

	printer, found := gw.State.Printers[msgData.PrinterID]
	if !found {
		res := &message.ResponseData{
			RequestID: msg.RequestID,
			Type:      msg.Type,
			Success:   false,
			Message:   "printer not found",
		}

		respMsg, err := message.New(message.TypeResponse, res)
		if err != nil {
			gw.logger.Error("failed to encode response message", "error", err)
			return
		}

		if err := gw.sendResponse(respMsg); err != nil {
			gw.logger.Error("failed to send response", "error", err)
			return
		}

		return
	}

	addr := fmt.Sprintf("%s:%s", printer.Addr, "9100")
	printerConn, err := dialPrinter(addr)
	if err != nil {
		res := &message.ResponseData{
			RequestID: msg.RequestID,
			Type:      msg.Type,
			Success:   false,
			Message:   fmt.Sprintf("failed to dial printer %s on %s", printer.Name, addr),
		}

		gw.logger.Error("failed to dial printer", "error", err, "response", res)

		respMsg, err := message.New(message.TypeResponse, res)
		if err != nil {
			gw.logger.Error("failed to encode response message", "error", err)
			return
		}

		if err := gw.sendResponse(respMsg); err != nil {
			gw.logger.Error("failed to send response", "error", err)
			return
		}

		return
	}

	if _, err := printerConn.Write(msg.Data); err != nil {
		res := &message.ResponseData{
			RequestID: msg.RequestID,
			Type:      msg.Type,
			Success:   false,
			Message:   fmt.Sprintf("failed to write to printer %s on %s", printer.Name, addr),
		}

		respMsg, err := message.New(message.TypeResponse, res)
		if err != nil {
			gw.logger.Error("failed to encode response message", "error", err)
			return
		}

		if err := gw.sendResponse(respMsg); err != nil {
			gw.logger.Error("failed to send response", "error", err)
			return
		}

		return
	}

	response := &message.ResponseData{
		RequestID: msg.RequestID,
		Type:      msg.Type,
		Message:   "message handled successfully",
		Success:   true,
	}

	respMsg, err := message.New(message.TypeResponse, response)
	if err != nil {
		gw.logger.Error("failed to encode response message", "error", err)
	}

	if err := gw.sendResponse(respMsg); err != nil {
		gw.logger.Error("failed to send response", "error", err)
	}
}

func (gw *Gateway) connect(addr string) error {
	if err := gw.loadCfg(); err != nil {
		return err
	}

	rawConn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	gw.logger.Info("connected to server")

	gw.conn = conn.Wrap(rawConn)
	return gw.authenticate()
}

func (gw *Gateway) authenticate() error {
	gw.conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	defer gw.conn.SetReadDeadline(time.Time{})

	auth := message.AuthData{
		Name:     gw.State.Name,
		Key:      gw.State.Key,
		Printers: []message.Printer{},
	}

	for id, printer := range gw.State.Printers {
		printer.ID = id
		auth.Printers = append(auth.Printers, printer)
	}

	msg, err := message.New(message.TypeAuth, auth)
	if err != nil {
		return err
	}

	data, err := message.Encode(msg)
	if err != nil {
		return err
	}

	if _, err := gw.conn.Write(append(data, '\n')); err != nil {
		return err
	}

	return nil
}

func (gw *Gateway) sendResponse(msg *message.Message) error {
	bytes, err := message.Encode(msg)
	if err != nil {
		return err
	}

	if _, err := gw.conn.Write(append(bytes, '\n')); err != nil {
		return err
	}

	return nil
}

func dialPrinter(addr string) (*net.TCPConn, error) {
	tcp, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp4", nil, tcp)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
