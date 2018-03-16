package site

import (
	"strconv"
	"time"

	"github.com/strongdm/comply/internal/model"
)

func loadValues() map[string]interface{} {
	values := make(map[string]interface{})

	values["Title"] = "Acme Compliance Program"
	values["Procedures"] = []string{
		"Jump",
		"Sit",
		"Squat",
	}

	stats := make(map[string]string)
	values["Stats"] = stats

	data, err := model.ReadData()
	if err == nil {
		var total, open, oldestDays, openProcess, openAudit, totalAudit, satisfiedControls, totalControls int

		// TODO: where does this go?
		satisfied := make(map[string]bool)
		for _, n := range data.Narratives {
			for _, controlKeys := range n.Satisfies {
				for _, key := range controlKeys {
					satisfied[key] = true
				}
			}
		}
		for _, n := range data.Policies {
			for _, controlKeys := range n.Satisfies {
				for _, key := range controlKeys {
					satisfied[key] = true
				}
			}
		}
		for _, n := range data.Procedures {
			for _, controlKeys := range n.Satisfies {
				for _, key := range controlKeys {
					satisfied[key] = true
				}
			}
		}

		for _, std := range data.Standards {
			totalControls += len(std.Controls)
			for controlKey := range std.Controls {
				if _, ok := satisfied[controlKey]; ok {
					satisfiedControls++
				}
			}
		}

		for _, t := range data.Tickets {
			total++

			if t.Bool("audit") {
				totalAudit++
			}

			if t.State == model.Open {
				if t.Bool("process") {
					openProcess++
				}
				if t.Bool("audit") {
					openAudit++
				}
				if t.CreatedAt != nil {
					oldestDays = int(time.Since(*t.CreatedAt).Hours() / float64(24))
				}
				open++
			}
		}

		stats["SatisfiedControls"] = strconv.Itoa(satisfiedControls)
		stats["TotalControls"] = strconv.Itoa(totalControls)
		stats["OldestDays"] = strconv.Itoa(oldestDays)
		stats["Total"] = strconv.Itoa(total)
		stats["Open"] = strconv.Itoa(open)
		stats["TotalAudit"] = strconv.Itoa(totalAudit)
		stats["OpenAudit"] = strconv.Itoa(openAudit)
		stats["OpenProcess"] = strconv.Itoa(openProcess)
		stats["ClosedAudit"] = strconv.Itoa(totalAudit - openAudit)
	}

	values["Narratives"] = data.Narratives
	values["Policies"] = data.Policies
	values["Procedures"] = data.Procedures
	return values
}
