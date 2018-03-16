package model

type Data struct {
	Tickets    []*Ticket
	Policies   []*Policy
	Procedures []*Procedure
	Audits     []*Audit
}

type Revision struct {
	Date    string `yaml:"date"`
	Comment string `yaml:"comment"`
}

type Satisfaction map[string][]string
