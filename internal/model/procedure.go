package model

import "time"

type Procedure struct {
	Name string `yaml:"name"`
	Code string `yaml:"code"`
	Cron string `yaml:"cron"`

	Revisions      []Revision   `yaml:"majorRevisions"`
	Satisfies      Satisfaction `yaml:"satisfies"`
	FullPath       string
	OutputFilename string
	ModifiedAt     time.Time
	Body           string
}
