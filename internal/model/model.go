package model

type Data struct {
	Tickets    []*Ticket
	Policies   []*Policy
	Procedures []*Procedure
	Audits     []*Audit
}

type Ticket struct {
	ID   string
	Name string
}

type Policy struct {
	ID      string
	Name    string
	Acronym string
}

type Procedure struct {
	ID   string
	Name string
}

type Audit struct {
	ID   string
	Name string
}
