package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/anvidev/inkbyte/internal/event"
)

var (
	socketPort = ":9091"
)

type client struct {
	mu    sync.Mutex
	conn  net.Conn
	State state `json:"state"`
}

type state struct {
	ApiKey   string            `json:"api_key"`
	Name     string            `json:"name"`
	Printers map[string]string `json:"printers"`
}

type printer struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

func new() *client {
	return &client{}
}

func (c *client) run() {
	if err := c.loadCfg(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	conn, err := net.Dial("tcp", socketPort)
	if err != nil {
		log.Fatalf("failed to dial socket server: %v", err)
	}

	c.conn = conn

	defer conn.Close()

	log.Printf("connected to socket on port %s: %+v", socketPort, c.State)

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		b := scanner.Bytes()

		var evt event.Event

		if err := json.Unmarshal(b, &evt); err != nil {
			log.Printf("error unmarshaling action: %v", err)
			continue
		}

		log.Printf("received event from server - type: %s, payload: %+v", evt.Name, evt.Payload)

		switch evt.Name {
		case "test":
			var testPayload event.TestPayload
			if err := json.Unmarshal(evt.Payload, &testPayload); err != nil {
				log.Printf("failed to unmarshal payload for action type %s: %v", "test", err)
				continue
			}
			c.handleTestMessage(testPayload)
		case "print":
			var payload event.PrintPayload
			if err := json.Unmarshal(evt.Payload, &payload); err != nil {
				log.Printf("failed to unmarshal payload for action type %s: %v", "test", err)
				continue
			}
			if err := c.handlePrintMessage(payload); err != nil {
				log.Printf("failed to handle print message: %v", err)
				continue
			}
		default:
			log.Printf("unknown event occured: %v", evt)
			continue
		}

		log.Println("handling complete")
	}

	if err := scanner.Err(); err != nil {
		log.Printf("connection error: %v", err)
	}

	log.Printf("disconnected client from socket")
}

func (c *client) handlePrintMessage(payload event.PrintPayload) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, found := c.State.Printers[payload.IP]
	if !found {
		if _, err := c.conn.Write([]byte("PRINTER NOT FOUND\n")); err != nil {
			return err
		}
		return fmt.Errorf("printer not found")
	}

	ip := fmt.Sprintf("%s:%s", payload.IP, "9100")
	log.Printf("attempting to dial %s", ip)

	conn, err := newTCPConn(ip)
	if err != nil {
		if _, err := c.conn.Write([]byte("PRINTER CONNECTION FAILED\n")); err != nil {
			return err
		}
		return fmt.Errorf("printer connection failed")
	}
	defer conn.Close()

	_, err = conn.Write([]byte(payload.Data))
	if err != nil {
		log.Println("error sending data to printer:", err)
		if _, err := c.conn.Write([]byte("PRINTER WRITE FAILED\n")); err != nil {
			return err
		}
		return err
	} else {
		log.Println("data sent to printer successfully")
	}

	if _, err := c.conn.Write([]byte("PRINT OK\n")); err != nil {
		return err
	}
	return nil
}

func (c *client) handleTestMessage(payload event.TestPayload) {
	log.Printf("handling test message: %v", payload)
}

func newTCPConn(ip string) (*net.TCPConn, error) {
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
