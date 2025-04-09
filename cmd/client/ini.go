package main

import (
	"strings"

	"github.com/anvidev/blaekward/internal/event"
	"github.com/google/uuid"
	"gopkg.in/ini.v1"
)

func (c *client) loadCfg() error {
	cfg, err := ini.Load("gateway.ini")
	if err != nil {
		return err
	}

	state := state{}

	state.Name = cfg.Section("gateway").Key("name").String()
	state.Key = cfg.Section("gateway").Key("key").String()

	printerCfgs := cfg.ChildSections("printer")
	var printers = make(map[string]event.Printer)

	for _, printerCfg := range printerCfgs {
		id := uuid.New().String()
		ip := printerCfg.Key("addr").String()
		name := strings.Split(printerCfg.Name(), ".")[1]
		p := event.Printer{
			Name: name,
			Addr: ip,
		}
		printers[id] = p
	}

	c.State = state
	c.State.Printers = printers

	return nil
}
