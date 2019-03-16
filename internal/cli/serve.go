package cli

import (
	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/render"
	"github.com/urfave/cli"
)

var serveCommand = cli.Command{
	Name:  "serve",
	Usage: "live updating version of the build command",
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:        "port",
			Value:       4000,
			Destination: &render.ServePort,
		},
	},
	Action: serveAction,
	Before: beforeAll(pandocMustExist, cleanContainers),
}

func serveAction(c *cli.Context) error {
	err := render.Build("output", true)
	if err != nil {
		return errors.Wrap(err, "serve failed")
	}
	return nil
}
