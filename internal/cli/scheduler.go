package cli

import (
	"github.com/strongdm/comply/internal/ticket"
	"github.com/urfave/cli"
)

var schedulerCommand = cli.Command{
	Name:   "scheduler",
	Usage:  "create tickets based on procedure schedule",
	Action: schedulerAction,
	Before: beforeAll(projectMustExist, ticketingMustBeConfigured),
}

func schedulerAction(c *cli.Context) error {
	err := syncAction(c)
	if err != nil {
		return err
	}
	return ticket.TriggerScheduled()
}
