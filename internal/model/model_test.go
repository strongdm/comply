package model

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMarshal(t *testing.T) {
	d := Data{
		Tickets: []*Ticket{
			{
				ID: "t1",
			},
		},
		Audits: []*Audit{
			{
				ID: "a1",
			},
		},
		Procedures: []*Procedure{
			{
				ID: "pro1",
			},
		},
		Policies: []*Document{
			{
				Name: "pol1",
			},
		},
	}
	m, _ := json.Marshal(d)
	encoded := string(m)
	if !strings.Contains(encoded, "t1") ||
		!strings.Contains(encoded, "a1") ||
		!strings.Contains(encoded, "pro1") ||
		!strings.Contains(encoded, "pol1") {
		t.Error("identifier not found in marshalled string")
	}
}
