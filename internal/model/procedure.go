package model

import (
	"fmt"
	"time"
)

var defaultLabels = []string{"comply", "comply-procedure"}

type Procedure struct {
	Name         string   `yaml:"name"`
	ID           string   `yaml:"id"`
	Cron         string   `yaml:"cron"`
	CustomLabels []string `yaml:"labels"`

	Revisions      []Revision   `yaml:"majorRevisions"`
	Satisfies      Satisfaction `yaml:"satisfies"`
	FullPath       string
	OutputFilename string
	ModifiedAt     time.Time
	Body           string
}

func (p *Procedure) Labels() []string {
	return append(defaultLabels, p.CustomLabels...)
}

func (p *Procedure) NewTicket() *Ticket {
	return &Ticket{
		Name: p.Name,
		Body: fmt.Sprintf("%s\n\n\n---\nProcedure-ID: %s", p.Body, p.ID),
	}
}
