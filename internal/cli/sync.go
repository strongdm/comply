package cli

import (
	"github.com/strongdm/comply/internal/model"
	"github.com/urfave/cli"
)

var syncCommand = cli.Command{
	Name:      "sync",
	ShortName: "s",
	Usage:     "sync external systems to local data cache",
	Action:    syncAction,
	Before:    projectMustExist,
}

func syncAction(c *cli.Context) error {
	tp := model.GetPlugin(model.Github)
	tickets, err := tp.FindByTagName("comply")
	if err != nil {
		panic(err)
	}
	for _, t := range tickets {
		model.DB().Write("tickets", t.ID, t)
	}
	return nil
}
