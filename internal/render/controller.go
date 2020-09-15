package render

import (
	"fmt"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/config"
	"github.com/strongdm/comply/internal/model"
)

type project struct {
	OrganizationName string
	Name             string
}

type stats struct {
	CriteriaTotal     int
	CriteriaSatisfied int

	ProcedureTotal      int
	ProcedureOpen       int
	ProcedureOldestDays int

	AuditOpen   int
	AuditClosed int
	AuditTotal  int
}

type renderData struct {
	// duplicates Project.OrganizationName
	Name       	string
	Project    	*project
	Stats      	*stats
	Narratives 	[]*model.Document
	Policies   	[]*model.Document
	Procedures 	[]*model.Procedure
	Frameworks	[]*model.Framework
	Tickets    	[]*model.Ticket
	Criteria  	[]*criterion
	Links      	*model.TicketLinks
}

type criterion struct {
	Framework    string
	CriteriaKey  string
	Name        string
	Description string
	Satisfied   bool
	SatisfiedBy []string
}

func load() (*model.Data, *renderData, error) {
	modelData, err := model.ReadData()
	if err != nil {
		return nil, nil, err
	}

	cfg := config.Config()
	project := &project{
		OrganizationName: cfg.Name,
		Name:             fmt.Sprintf("%s Compliance Program", cfg.Name),
	}

	satisfied := model.CriteriaSatisfied(modelData)
	criteria := make([]*criterion, 0)
	for _, framework := range modelData.Frameworks {
		for key, c := range framework.Criteria{
			satisfactions, ok := satisfied[key]
			satisfied := ok && len(satisfactions) > 0
			criteria = append(criteria, &criterion{
				Framework:    framework.Name,
				CriteriaKey:  key,
				Name:        c.Name,
				Description: c.Description,
				Satisfied:   satisfied,
				SatisfiedBy: satisfactions,
			})
		}
	}
	sort.Slice(criteria, func(i, j int) bool {
		return criteria[i].CriteriaKey < criteria[j].CriteriaKey
	})

	rd := &renderData{}
	rd.Narratives = modelData.Narratives
	rd.Policies = modelData.Policies
	rd.Procedures = modelData.Procedures
	rd.Frameworks = modelData.Frameworks
	rd.Tickets = modelData.Tickets
	rd.Links = &model.TicketLinks{}
	rd.Project = project
	rd.Name = project.OrganizationName
	rd.Criteria= criteria

	ts, err := config.Config().TicketSystem()
	if err != nil {
		return nil, nil, errors.Wrap(err, "error in ticket system configuration")
	}

	tp := model.GetPlugin(model.TicketSystem(ts))
	if tp.Configured() {
		links := tp.Links()
		rd.Links = &links
	}

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

	satisfied := model.CriteriaSatisfied(modelData)

	for _, std := range renderData.Frameworks {
		stats.CriteriaTotal += len(std.Criteria)
		for criteriaKey := range std.Criteria{
			if _, ok := satisfied[criteriaKey]; ok {
				stats.CriteriaSatisfied++
			}
		}
	}

	for _, t := range renderData.Tickets {
		if t.Bool("audit") {
			stats.AuditTotal++
		}

		if t.State == model.Open {
			if t.Bool("comply-procedure") {
				stats.ProcedureOpen++
				if t.CreatedAt != nil {
					age := int(time.Since(*t.CreatedAt).Hours() / float64(24))
					if stats.ProcedureOldestDays < age {
						stats.ProcedureOldestDays = age
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
