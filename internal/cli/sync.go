package cli

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/strongdm/comply/internal/model"
	"github.com/urfave/cli"
)

var syncCommand = cli.Command{
	Name:      "sync",
	ShortName: "s",
	Usage:     "sync external systems to local data cache",
	Action:    syncAction,
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

	policies := model.ReadPolicies()
	fmt.Println("LEN ", len(policies))
	for _, p := range policies {
		spew.Dump(p)
	}

	return nil
}
