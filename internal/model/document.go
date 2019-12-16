package model

import "time"

type Document struct {
	Name    string `yaml:"name"`
	Acronym string `yaml:"acronym"`

	Revisions      []Revision   `yaml:"majorRevisions"`
	Satisfies      Satisfaction `yaml:"satisfies"`
	Live           bool         `yaml:"live"`
	FullPath       string
	OutputFilename string
	ModifiedAt     time.Time
	Body           string
}
