package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"

	"github.com/anvidev/inkbyte/internal/event"
)

var (
	socketPort = ":9091"
)

type client struct {
	conn net.Conn
}

func new() *client {
	return &client{}
}

func (c *client) run() {
	conn, err := net.Dial("tcp", socketPort)
	if err != nil {
		log.Fatalf("failed to dial socket server: %v", err)
	}

	c.conn = conn

	defer conn.Close()

	log.Printf("connected to socket on port %s", socketPort)

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
	log.Printf("handling test message: %v", payload)
	if _, err := c.conn.Write([]byte("PRINT OK\n")); err != nil {
		return err
	}
	return nil
}

func (c *client) handleTestMessage(payload event.TestPayload) {
	log.Printf("handling test message: %v", payload)
}
