package ticket

import (
	"fmt"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/strongdm/comply/internal/config"
	"github.com/strongdm/comply/internal/model"
)

func byProcedureByTime(tickets []*model.Ticket) map[string][]*model.Ticket {
	result := make(map[string][]*model.Ticket)
	for _, t := range tickets {
		procedureID := t.ProcedureID()
		if procedureID == "" {
			// missing procedure metadata; skip
			continue
		}
		list, ok := result[procedureID]
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
		result[procedureID] = list
	}
	return result
}

func TriggerScheduled() error {
	rawTickets, err := model.ReadTickets()
	if err != nil {
		return err
	}
	tickets := byProcedureByTime(rawTickets)
	procedures, err := model.ReadProcedures()
	if err != nil {
		return err
	}

	for _, procedure := range procedures {
		if procedure.Cron == "" {
			continue
		}

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
			err = trigger(procedure)
			if err != nil {
				return err
			}
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
				err = trigger(procedure)
				if err != nil {
					return err
				}
				break SEARCH
			}
		}

	}
	return nil
}

func trigger(procedure *model.Procedure) error {
	fmt.Printf("triggering procedure %s (cron expression: %s)\n", procedure.Name, procedure.Cron)

	ts, err := config.Config().TicketSystem()
	if err != nil {
		return errors.Wrap(err, "error in ticket system configuration")
	}

	tp := model.GetPlugin(model.TicketSystem(ts))
	err = tp.Create(procedure.NewTicket(), procedure.Labels())
	return err
}
