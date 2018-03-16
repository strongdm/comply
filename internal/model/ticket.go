package model

import "time"

type TicketState string

const (
	Open   = TicketState("open")
	Closed = TicketState("closed")
)

type Ticket struct {
	ID         string
	Name       string
	State      TicketState
	Body       string
	Attributes map[string]interface{}
	ClosedAt   *time.Time
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

func (t *Ticket) SetBool(name string) {
	t.Attributes[name] = true
}
func (t *Ticket) Bool(name string) bool {
	bi, ok := t.Attributes[name]
	if !ok {
		return false
	}

	b, ok := bi.(bool)
	if !ok {
		return false
	}

	return b
}
