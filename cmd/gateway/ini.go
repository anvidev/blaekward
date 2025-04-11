package main

import (
	"strings"

	"github.com/anvidev/blaekward/internal/id"
	"github.com/anvidev/blaekward/internal/message"
	"gopkg.in/ini.v1"
)

func (gw *Gateway) loadCfg() error {
	cfg, err := ini.Load("gateway.ini")
	if err != nil {
		return err
	}

	state := state{}

	state.Name = cfg.Section("gateway").Key("name").String()
	state.Key = cfg.Section("gateway").Key("key").String()

	printerCfgs := cfg.ChildSections("printer")
	var printers = make(map[string]message.Printer)

	for _, printerCfg := range printerCfgs {
		id := id.Generate(8, id.Numbers)
		ip := printerCfg.Key("addr").String()
		name := strings.Split(printerCfg.Name(), ".")[1]
		p := message.Printer{
			Name: name,
			Addr: ip,
		}
		printers[id] = p
	}

	gw.State = state
	gw.State.Printers = printers

	return nil
}
