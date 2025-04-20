package main

import (
	"strings"

	"github.com/anvidev/blaekward/internal/event"
	"github.com/anvidev/blaekward/internal/id"
	"gopkg.in/ini.v1"
)

func (gw *Gateway) loadIniFile() error {
	cfg, err := ini.Load("gateway.ini")
	if err != nil {
		return err
	}

	state := setup{}

	state.Name = cfg.Section("gateway").Key("name").String()
	state.Key = cfg.Section("gateway").Key("key").String()

	printerCfgs := cfg.ChildSections("printer")
	var printers = make(map[string]event.Printer)

	for _, printerCfg := range printerCfgs {
		id := id.String(8, id.Numbers)
		ip := printerCfg.Key("addr").String()
		name := strings.Split(printerCfg.Name(), ".")[1]
		p := event.Printer{
			Name: name,
			Addr: ip,
		}
		printers[id] = p
	}

	gw.setup = state
	gw.setup.Printers = printers

	return nil
}
