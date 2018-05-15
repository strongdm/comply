package cli

import (
	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/render"
	"github.com/urfave/cli"
)

var serveCommand = cli.Command{
	Name:   "serve",
	Usage:  "live updating version of the build command",
	Action: serveAction,
	Before: beforeAll(dockerMustExist, cleanContainers),
}

func serveAction(c *cli.Context) error {
	err := render.Build("output", true)
	if err != nil {
		return errors.Wrap(err, "serve failed")
	}
	return nil
}
