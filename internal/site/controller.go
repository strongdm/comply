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

	rt, err := model.DB().ReadAll("tickets")
	if err == nil {
		ts := model.Tickets(rt)
		var total, open, oldestDays int
		for _, t := range ts {
			total++
			if t.State == model.Open {
				if t.CreatedAt != nil {
					oldestDays = int(time.Since(*t.CreatedAt).Hours() / float64(24))
				}
				open++
			}

		}

		values["OldestDays"] = strconv.Itoa(oldestDays)
		values["Total"] = strconv.Itoa(total)
		values["Open"] = strconv.Itoa(open)
	}

	policies := model.ReadPolicies()

	values["Policies"] = policies

	return values
}
