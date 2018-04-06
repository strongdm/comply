package ticket

import (
	"fmt"
	"sort"
	"time"

	"github.com/robfig/cron"
	"github.com/strongdm/comply/internal/model"
)

func byProcessByTime(tickets []*model.Ticket) map[string][]*model.Ticket {
	result := make(map[string][]*model.Ticket)
	for _, t := range tickets {
		processID := t.ProcedureID()
		if processID == "" {
			// missing process metadata; skip
			continue
		}
		list, ok := result[processID]
		if !ok {
			list = make([]*model.Ticket, 0)
		}
		list = append(list, t)
		sort.Slice(list, func(i, j int) bool {
			if list[i].CreatedAt == nil || list[j].CreatedAt == nil {
				return false
			}
			return list[i].CreatedAt.Before(*list[j].CreatedAt)
		})
		result[processID] = list
	}
	return result
}

func TriggerScheduled() {
	rawTickets := model.ReadTickets()
	tickets := byProcessByTime(rawTickets)
	// spew.Dump(tickets)
	for _, procedure := range model.ReadProcedures() {
		if procedure.Cron == "" {
			continue
		}

		// TODO
		procedureID := procedure.ID
		schedule, err := cron.Parse(procedure.Cron)
		if err != nil {
			continue
		}
		ticketsForProc, ok := tickets[procedureID]
		if ok {
			// find most recent one
			mostRecent := ticketsForProc[len(ticketsForProc)-1]
			if mostRecent.CreatedAt == nil {
				continue
			}

			// would another have triggered since?
			nextTrigger := schedule.Next(*mostRecent.CreatedAt).UTC()
			if nextTrigger.After(time.Now().UTC()) {
				// in the future, nothing to do
				continue
			}
			trigger(procedure)
		} else {
			// don't go back further than 13 months
			tooOld := time.Now().Add(-1 * time.Hour * 24 * (365 + 30))
			// search back one day until triggers
			triggeredAt := time.Now().Add(-24 * time.Hour).UTC()
		SEARCH:
			for {
				if triggeredAt.Before(tooOld) {
					break SEARCH
				}

				candidate := schedule.Next(triggeredAt)
				// in the future? not far eonugh back yet.
				if candidate.After(time.Now().UTC()) {
					triggeredAt = triggeredAt.Add(-24 * time.Hour)
					continue
				}

				// is in the past? then trigger.
				trigger(procedure)
				break SEARCH
			}
		}

	}
	return
}

func trigger(procedure *model.Procedure) {
	// TODO: don't hardcode GH
	tp := model.GetPlugin(model.Github)
	tp.Create(&model.Ticket{
		Name: procedure.Name,
		Body: fmt.Sprintf("%s\n\n\n---\nProcedure-ID: %s", procedure.Body, procedure.ID),
	}, []string{"comply", "comply-procedure"})
}
