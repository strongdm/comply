package render

import (
	"fmt"
	"sort"
	"time"

	"github.com/strongdm/comply/internal/config"
	"github.com/strongdm/comply/internal/model"
)

type project struct {
	OrganizationName string
	Name             string
}

type stats struct {
	ControlsTotal     int
	ControlsSatisfied int

	ProcessTotal      int
	ProcessOpen       int
	ProcessOldestDays int

	AuditOpen   int
	AuditClosed int
	AuditTotal  int
}

type renderData struct {
	Project    *project
	Stats      *stats
	Narratives []*model.Narrative
	Policies   []*model.Policy
	Procedures []*model.Procedure
	Standards  []*model.Standard
	Tickets    []*model.Ticket
	Controls   []*control
}

type control struct {
	Standard    string
	ControlKey  string
	Name        string
	Description string
	Satisfied   bool
}

func load() (*renderData, error) {
	modelData := model.ReadData()

	cfg := config.Config()
	project := &project{
		OrganizationName: cfg.Name,
		Name:             fmt.Sprintf("%s Compliance Program", cfg.Name),
	}

	satisfied := model.ControlsSatisfied(modelData)
	controls := make([]*control, 0)
	for _, standard := range modelData.Standards {
		for key, c := range standard.Controls {
			controls = append(controls, &control{
				Standard:    standard.Name,
				ControlKey:  key,
				Name:        c.Name,
				Description: c.Description,
				Satisfied:   satisfied[key],
			})
		}
	}
	sort.Slice(controls, func(i, j int) bool {
		return controls[i].ControlKey < controls[j].ControlKey
	})

	rd := &renderData{}
	rd.Narratives = modelData.Narratives
	rd.Policies = modelData.Policies
	rd.Procedures = modelData.Procedures
	rd.Standards = modelData.Standards
	rd.Tickets = modelData.Tickets
	rd.Project = project
	rd.Controls = controls
	return rd, nil
}

func loadWithStats() (*renderData, error) {
	d, err := load()
	if err != nil {
		return nil, err
	}

	addStats(d)
	return d, nil
}

func addStats(data *renderData) {
	stats := &stats{}

	satisfied := make(map[string]bool)
	for _, n := range data.Narratives {
		for _, controlKeys := range n.Satisfies {
			for _, key := range controlKeys {
				satisfied[key] = true
			}
		}
	}
	for _, n := range data.Policies {
		for _, controlKeys := range n.Satisfies {
			for _, key := range controlKeys {
				satisfied[key] = true
			}
		}
	}
	for _, n := range data.Procedures {
		for _, controlKeys := range n.Satisfies {
			for _, key := range controlKeys {
				satisfied[key] = true
			}
		}
	}

	for _, std := range data.Standards {
		stats.ControlsTotal += len(std.Controls)
		for controlKey := range std.Controls {
			if _, ok := satisfied[controlKey]; ok {
				stats.ControlsSatisfied++
			}
		}
	}

	for _, t := range data.Tickets {
		if t.Bool("audit") {
			stats.AuditTotal++
		}

		if t.State == model.Open {
			if t.Bool("process") {
				stats.ProcessOpen++
				if t.CreatedAt != nil {
					age := int(time.Since(*t.CreatedAt).Hours() / float64(24))
					if stats.ProcessOldestDays < age {
						stats.ProcessOldestDays = age
					}
				}
			}
			if t.Bool("audit") {
				stats.AuditOpen++
			}
		}
	}

	data.Stats = stats
}
