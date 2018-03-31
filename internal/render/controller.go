package render

import (
	"fmt"
	"time"

	"github.com/strongdm/comply/internal/config"
	"github.com/strongdm/comply/internal/model"
)

// Project contains all project-scope fields
type Project struct {
	OrganizationName string
	Name             string
}

// Stats contains all computed counts
type Stats struct {
	ControlsTotal     int
	ControlsSatisfied int

	ProcessTotal      int
	ProcessOpen       int
	ProcessOldestDays int

	AuditOpen   int
	AuditClosed int
	AuditTotal  int
}

func loadValues() map[string]interface{} {
	stats := &Stats{}
	data, err := model.ReadData()
	if err == nil {
		// TODO: where does this go?
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
	}

	cfg := config.Config()
	project := Project{
		OrganizationName: cfg.Name,
		Name:             fmt.Sprintf("%s Compliance Program", cfg.Name),
	}

	values := make(map[string]interface{})
	values["Project"] = project
	values["Stats"] = stats
	values["Narratives"] = data.Narratives
	values["Policies"] = data.Policies
	values["Procedures"] = data.Procedures
	return values
}
