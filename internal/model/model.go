package model

type Data struct {
	Frameworks	[]*Framework
	Narratives 	[]*Document
	Policies   	[]*Document
	Procedures 	[]*Procedure
	Tickets    	[]*Ticket
	Audits     	[]*Audit
}

type Revision struct {
	Date    string `yaml:"date"`
	Comment string `yaml:"comment"`
}

type Satisfaction map[string][]string
