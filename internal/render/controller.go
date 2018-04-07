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

func load() (*model.Data, *renderData, error) {
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
	return modelData, rd, nil
}

func loadWithStats() (*model.Data, *renderData, error) {
	modelData, renderData, err := load()
	if err != nil {
		return nil, nil, err
	}

	addStats(modelData, renderData)
	return modelData, renderData, nil
}

func addStats(modelData *model.Data, renderData *renderData) {
	stats := &stats{}

	satisfied := model.ControlsSatisfied(modelData)

	for _, std := range renderData.Standards {
		stats.ControlsTotal += len(std.Controls)
		for controlKey := range std.Controls {
			if _, ok := satisfied[controlKey]; ok {
				stats.ControlsSatisfied++
			}
		}
	}

	for _, t := range renderData.Tickets {
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

	renderData.Stats = stats
}
