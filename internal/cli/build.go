package cli

import (
	"github.com/strongdm/comply/internal/render"
	"github.com/urfave/cli"
)

var buildCommand = cli.Command{
	Name:      "build",
	ShortName: "b",
	Usage:     "generate a static website summarizing the compliance program",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "live, l",
			Usage: "rebuild static site after filesystem changes",
		},
	},
	Action: buildAction,
	Before: beforeAll(projectMustExist, dockerMustExist),
}

func buildAction(c *cli.Context) error {
	err := render.Build("output", false)
	if err != nil {
		panic(err)
	}
	return nil
}
