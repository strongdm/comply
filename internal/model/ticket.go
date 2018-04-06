package model

import (
	"strings"
	"time"
)

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

func (t *Ticket) ProcedureID() string {
	md := t.metadata()
	if v, ok := md["Procedure-ID"]; ok {
		return v
	}
	return ""
}

func (t *Ticket) metadata() map[string]string {
	md := make(map[string]string)
	lines := strings.Split(t.Body, "\n")
	for _, line := range lines {
		// TODO: transition to RFC822 parsing
		if strings.Contains(line, ":") {
			tokens := strings.Split(line, ":")
			if len(tokens) != 2 {
				continue
			}
			md[strings.TrimSpace(tokens[0])] = strings.TrimSpace(tokens[1])
		}
	}
	return md
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
