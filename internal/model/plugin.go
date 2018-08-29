package model

import (
	"fmt"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/strongdm/comply/internal/config"
)

var tsPluginsMu sync.Mutex
var tsPlugins = make(map[TicketSystem]TicketPlugin)
var tsConfigureOnce sync.Once

// TicketSystem is the type of ticket database.
type TicketSystem string

const (
	// Jira from Atlassian.
	Jira = TicketSystem(config.Jira)
	// GitHub from GitHub.
	GitHub = TicketSystem(config.GitHub)
	// GitLab from GitLab.
	GitLab = TicketSystem(config.GitLab)
	// NoTickets indicates no ticketing system integration.
	NoTickets = TicketSystem(config.NoTickets)
)

type TicketLinks struct {
	ProcedureOpen string
	ProcedureAll  string
	AuditOpen     string
	AuditAll      string
}

// TicketPlugin models support for ticketing systems.
type TicketPlugin interface {
	Get(ID string) (*Ticket, error)
	FindOpen() ([]*Ticket, error)
	FindByTag(name, value string) ([]*Ticket, error)
	FindByTagName(name string) ([]*Ticket, error)
	Create(ticket *Ticket, labels []string) error
	Configure(map[string]interface{}) error
	Prompts() map[string]string
	Links() TicketLinks
	LinkFor(ticket *Ticket) string
	Configured() bool
}

// GetPlugin loads the ticketing database.
func GetPlugin(ts TicketSystem) TicketPlugin {
	tsPluginsMu.Lock()
	defer tsPluginsMu.Unlock()

	if ts == NoTickets {
		return &noopTicketSystem{}
	}

	tp, ok := tsPlugins[ts]
	if !ok {
		panic("Unknown ticket system: " + ts)
	}

	if config.Exists() {
		tsConfigureOnce.Do(func() {
			ticketsMap := config.Config().Tickets
			hasTickets := true

			cfg, ok := ticketsMap[string(ts)]
			if !ok {
				hasTickets = false
			}

			if hasTickets {
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
				err := tp.Configure(cfgStringed)
				if err != nil {
					panic(fmt.Sprintf("Configuration error `%s` in project YAML", err))
				}
			}
		})
	}

	return tp
}

// Register ticketing system plugin.
func Register(ts TicketSystem, plugin TicketPlugin) {
	tsPluginsMu.Lock()
	defer tsPluginsMu.Unlock()
	_, ok := tsPlugins[ts]
	if ok {
		panic("Duplicate ticketing system registration: " + ts)
	}

	tsPlugins[ts] = plugin
}

type noopTicketSystem struct{}

func (*noopTicketSystem) Get(ID string) (*Ticket, error) {
	return nil, nil
}
func (*noopTicketSystem) FindOpen() ([]*Ticket, error) {
	return []*Ticket{}, nil
}
func (*noopTicketSystem) FindByTag(name, value string) ([]*Ticket, error) {
	return []*Ticket{}, nil
}
func (*noopTicketSystem) FindByTagName(name string) ([]*Ticket, error) {
	return []*Ticket{}, nil
}
func (*noopTicketSystem) Create(ticket *Ticket, labels []string) error {
	return nil
}
func (*noopTicketSystem) Configure(map[string]interface{}) error {
	return nil
}
func (*noopTicketSystem) Prompts() map[string]string {
	return make(map[string]string)
}
func (*noopTicketSystem) Links() TicketLinks {
	return TicketLinks{}
}
func (*noopTicketSystem) LinkFor(ticket *Ticket) string {
	return ""
}
func (*noopTicketSystem) Configured() bool {
	return false
}
