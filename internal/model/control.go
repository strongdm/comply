package model

import (
  "time"
	"html/template"
)

type Control struct {
	Name            string      `yaml:"name"`
	ID              string      `yaml:"identifier"`
  Family          string      `yaml:"family"`
  Owner           string      `yaml:"owner"`
  GoverningPolicy []Policy    `yaml:"governingPolicy"`
  Revisions       []Revision  `yaml:"revisions"`
  Targets         Target      `yaml:"targets"`
  Published       string      `yaml:"published"`

  FullPath       string
	OutputFilename string
	ModifiedAt     time.Time

	Body            string
  BodyHTML        template.HTML
}
