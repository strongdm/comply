package cli

import (
	"github.com/strongdm/comply/internal/site"
	"github.com/urfave/cli"
)

var buildCommand = cli.Command{
	Name:  "build",
	Usage: "generate a static website summarizing the compliance program",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "live, l",
			Usage: "rebuild static site after filesystem changes",
		},
	},
	Action: buildAction,
}

func buildAction(c *cli.Context) error {
	live := false
	input := ""
	if c.Bool("live") {
		live = true
		input = "site"
	}
	// create directory structure
	err := site.Build(input, "output", live)
	if err != nil {
		panic(err)
	}
	return nil
}
