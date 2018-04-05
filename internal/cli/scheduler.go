package cli

import (
	"github.com/strongdm/comply/internal/ticket"
	"github.com/urfave/cli"
)

var schedulerCommand = cli.Command{
	Name:   "scheduler",
	Usage:  "create tickets based on procedure schedule",
	Action: schedulerAction,
	Before: projectMustExist,
}

func schedulerAction(c *cli.Context) error {
	ticket.TriggerScheduled()
	return nil
}
