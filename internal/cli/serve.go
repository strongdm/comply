package cli

import (
	"github.com/strongdm/comply/internal/render"
	"github.com/urfave/cli"
)

var serveCommand = cli.Command{
	Name:   "serve",
	Usage:  "live updating version of the build command",
	Action: serveAction,
	Before: beforeAll(projectMustExist, dockerMustExist),
}

func serveAction(c *cli.Context) error {
	err := render.Build("output", true)
	if err != nil {
		panic(err)
	}
	return nil
}
