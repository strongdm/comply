package model

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMarshal(t *testing.T) {
	d := Data{
		Tickets: []*Ticket{
			&Ticket{
				ID: "t1",
			},
		},
		Audits: []*Audit{
			&Audit{
				ID: "a1",
			},
		},
		Procedures: []*Procedure{
			&Procedure{
				Code: "pro1",
			},
		},
		Policies: []*Policy{
			&Policy{
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
