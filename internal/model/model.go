package model

import "time"

type Data struct {
	Tickets    []*Ticket
	Policies   []*Policy
	Procedures []*Procedure
	Audits     []*Audit
}

type TicketState string

const (
	Open   = TicketState("open")
	Closed = TicketState("closed")
)

type Ticket struct {
	ID        string
	Name      string
	State     TicketState
	Body      string
	ClosedAt  *time.Time
	CreatedAt *time.Time
	UpdatedAt *time.Time
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
