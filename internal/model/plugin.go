package model

import (
	"fmt"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/strongdm/comply/internal/config"
)

var tsPluginsMu sync.Mutex
var tsPlugins = make(map[TicketSystem]TicketPlugin)

// TicketSystem is the type of ticket database
type TicketSystem string

const (
	// JIRA from Atlassian
	JIRA = TicketSystem("jira")
	// Github from Github
	Github = TicketSystem("github")
)

// TicketPlugin models support for ticketing systems
type TicketPlugin interface {
	Get(ID string) (*Ticket, error)
	FindOpen() ([]*Ticket, error)
	FindByTag(name, value string) ([]*Ticket, error)
	FindByTagName(name string) ([]*Ticket, error)
	Create(*Ticket) error
	Update(*Ticket) error
	Close(ID string) error
	Configure(map[string]interface{}) error
}

// GetPlugin loads the ticketing database
func GetPlugin(ts TicketSystem) TicketPlugin {
	tsPluginsMu.Lock()
	defer tsPluginsMu.Unlock()

	tp, ok := tsPlugins[ts]
	if !ok {
		panic("Unknown ticket system: " + ts)
	}
	return tp
}

// Register ticketing system plugin
func Register(ts TicketSystem, plugin TicketPlugin) {
	tsPluginsMu.Lock()
	defer tsPluginsMu.Unlock()
	_, ok := tsPlugins[ts]
	if ok {
		panic("Duplicate ticketing system registration: " + ts)
	}

	// TODO: move parsing and feedback to YAML package
	yml := config.YAML()
	tickets, ok := yml["tickets"]
	if !ok {
		spew.Dump(yml)
		panic("Missing ticket configuration; add `tickets` block to project YAML")
	}

	ticketsMap, ok := tickets.(map[interface{}]interface{})
	if !ok {
		spew.Dump(tickets)
		panic("Malformed ticket configuration; modify `tickets` block in project YAML")
	}

	cfg, ok := ticketsMap[string(ts)]
	if !ok {
		spew.Dump(cfg)
		panic(fmt.Sprintf("Missing configuration for plugin system; add `%s` block to project YAML", string(ts)))
	}

	cfgTyped, ok := cfg.(map[interface{}]interface{})
	if !ok {
		spew.Dump(cfg)
		panic(fmt.Sprintf("malformatted ticket configuration block `%s` in project YAML", string(ts)))
	}

	cfgStringed := make(map[string]interface{})
	for k, v := range cfgTyped {
		kS, ok := k.(string)
		if !ok {
			spew.Dump(cfgStringed)
			panic(fmt.Sprintf("malformatted key in configuration block `%s` in project YAML", string(ts)))
		}
		cfgStringed[kS] = v
	}

	plugin.Configure(cfgStringed)

	tsPlugins[ts] = plugin
}
