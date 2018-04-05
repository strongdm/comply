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
		processID := t.ProcessID()
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
		procedureID := "Fish"
		schedule, err := cron.Parse(procedure.Cron)
		if err != nil {
			continue
		}
		ticketsForProc, ok := tickets[procedureID]
		if ok {
			fmt.Println("OK")
			// find most recent one
			mostRecent := ticketsForProc[len(ticketsForProc)-1]
			fmt.Println("MR IS ", mostRecent.CreatedAt)
			if mostRecent.CreatedAt == nil {
				continue
			}

			// would another have triggered since?
			nextTrigger := schedule.Next(*mostRecent.CreatedAt).UTC()
			fmt.Println("NT IS ", nextTrigger)
			if nextTrigger.After(time.Now().UTC()) {
				// in the future, nothing to do
				continue
			}
			trigger(procedure)
		} else {
			// walk back one day until triggers
			fmt.Println("NAW BRA")

		}

	}
	return
}

func trigger(procedure *model.Procedure) {
	fmt.Println("TRIGGERING " + procedure.Name)
}
