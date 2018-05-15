package model

import "time"

type Narrative struct {
	Name    string `yaml:"name"`
	Acronym string `yaml:"acronym"`

	Revisions      []Revision   `yaml:"majorRevisions"`
	Satisfies      Satisfaction `yaml:"satisfies"`
	FullPath       string
	OutputFilename string
	ModifiedAt     time.Time
	Body           string
}
