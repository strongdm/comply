package cli

import (
	"github.com/strongdm/comply/internal/config"
	"github.com/strongdm/comply/internal/model"
	"github.com/urfave/cli"
)

var syncCommand = cli.Command{
	Name:   "sync",
	Usage:  "sync ticket status to local cache",
	Action: syncAction,
	Before: beforeAll(projectMustExist, ticketingMustBeConfigured),
}

func syncAction(c *cli.Context) error {
	ts, err := config.Config().TicketSystem()
	if err != nil {
		return cli.NewExitError("error in ticket system configuration", 1)
	}

	tp := model.GetPlugin(model.TicketSystem(ts))
	tickets, err := tp.FindByTagName("comply")
	if err != nil {
		return err
	}
	for _, t := range tickets {
		err = model.DB().Write("tickets", t.ID, t)
		if err != nil {
			return err
		}
	}
	return nil
}
