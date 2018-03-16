package model

type Narrative struct {
	Name      string       `yaml:"name"`
	Revisions []Revision   `yaml:"majorRevisions"`
	Satisfies Satisfaction `yaml:"satisfies"`
}
