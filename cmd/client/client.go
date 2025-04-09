package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/anvidev/blaekward/internal/event"
	"github.com/anvidev/blaekward/internal/utils"
	"github.com/anvidev/blaekward/internal/zpl"
)

var (
	socketPort = ":9091"
)

type client struct {
	logger *slog.Logger

	mu    sync.Mutex
	conn  net.Conn
	State state `json:"state"`
}

type state struct {
	Key      string                   `json:"key"`
	CID      string                   `json:"cid"`
	Name     string                   `json:"name"`
	Printers map[string]event.Printer `json:"printers"`
}

func new() *client {
	return &client{logger: slog.Default()}
}

func (c *client) run() {
	if err := c.loadCfg(); err != nil {
		c.logger.Error("failed to load config", "error", err)
		return
	}

	url := "blaekward.anvi.dev"

	conn, err := net.Dial("tcp", url+socketPort)
	if err != nil {
		c.logger.Error("failed to dial socket server", "error", err)
		return
	}

	c.conn = conn

	defer conn.Close()

	if err := c.authenticate(); err != nil {
		c.logger.Error("failed to authenticate", "error", err)
		return
	}

	c.logger.Info("connected to socket", "port", socketPort, "state", c.State)

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		b := scanner.Bytes()

		var evt event.Event

		if err := json.Unmarshal(b, &evt); err != nil {
			c.logger.Error("error unmarshaling event", "error", err)
			continue
		}

		c.logger.Info("received event from server", "event", evt.Event)

		switch evt.Event {
		case "print":
			var payload event.PrintPayload
			if err := json.Unmarshal(evt.Payload, &payload); err != nil {
				c.logger.Error("failed to unmarshal event payload", "event", evt.Event, "error", err)
				continue
			}
			if err := c.handlePrintMessage(payload); err != nil {
				c.logger.Error("event failed", "event", evt.Event, "error", err)
				continue
			}
		default:
			c.logger.Warn("unknown event occured", "event", evt.Event)
			continue
		}

		c.logger.Info("event handling complete", "event", evt.Event)
	}

	if err := scanner.Err(); err != nil {
		c.logger.Error("connection error", "error", err)
	}

	c.logger.Info("disconnected client from socket")
}

func (c *client) handlePrintMessage(payload event.PrintPayload) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	printer, found := c.State.Printers[payload.PrinterID]
	if !found {
		c.logger.Warn("printer not found", "addr", payload.PrinterID)
		return c.sendSocketResponse(utils.Response{
			Success: false,
			Message: fmt.Sprintf("printer %s not found", payload.PrinterID),
		})
	}

	addr := fmt.Sprintf("%s:%s", printer.Addr, "9100")

	if err := c.sendPrinterJob(addr, payload.Data); err != nil {
		return c.sendSocketResponse(utils.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.sendSocketResponse(utils.Response{
		Success: true,
		Message: "Print successful",
	})
}

func connectToPrinter(ip string) (*net.TCPConn, error) {
	tcp, err := net.ResolveTCPAddr("tcp4", ip)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp4", nil, tcp)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *client) sendPrinterJob(addr, data string) error {
	conn, err := connectToPrinter(addr)
	if err != nil {
		return err
	}

	z := zpl.New(76, 51)
	z.CF(
		zpl.WithCFFont("A"),
		zpl.WithCFCharHeight(30),
	)

	z.FO(30, 30)
	z.FD("John Doe")
	z.FS()

	z.FO(30, 65)
	z.FD("100 Main Street")
	z.FS()

	z.FO(30, 100)
	z.FD("Springfield TN 39021")
	z.FS()

	z.FO(30, 135)
	z.FD("United States")
	z.FS()

	z.BY(
		zpl.WithModuleWidth(3),
		zpl.WithBarcodeHeight(135),
	)
	z.FO(30, 200)
	z.BC()
	z.FD("12345678")

	bytes := z.Bytes()

	defer conn.Close()

	_, err = conn.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) sendSocketResponse(r utils.Response) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}

	if _, err := c.conn.Write(append(data, '\n')); err != nil {
		return err
	}

	return nil
}

func (c *client) authenticate() error {
	c.conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	defer c.conn.SetReadDeadline(time.Time{})

	authEvent := event.AuthenticationEvent{
		Name:     c.State.Name,
		CID:      c.State.CID,
		Key:      c.State.Key,
		Printers: []event.Printer{},
	}

	for id, printer := range c.State.Printers {
		printer.ID = id
		authEvent.Printers = append(authEvent.Printers, printer)
	}

	data, err := authEvent.Encode()
	if err != nil {
		return err
	}

	if _, err := c.conn.Write([]byte(data)); err != nil {
		return err
	}

	return nil
}
