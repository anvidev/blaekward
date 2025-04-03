package main

import (
	"strings"

	"gopkg.in/ini.v1"
)

func (c *client) loadCfg() error {
	cfg, err := ini.Load("gateway.ini")
	if err != nil {
		return err
	}

	state := state{}

	state.Name = cfg.Section("gateway").Key("name").String()
	state.ApiKey = cfg.Section("gateway").Key("apikey").String()

	printerCfgs := cfg.ChildSections("printer")
	var printers = make(map[string]string)

	for _, printerCfg := range printerCfgs {
		ip := printerCfg.Key("ip").String()
		name := strings.Split(printerCfg.Name(), ".")[1]
		printers[ip] = name
	}

	c.State = state
	c.State.Printers = printers

	return nil
}
