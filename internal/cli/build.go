package cli

import (
	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/render"
	"github.com/urfave/cli"
)

var buildCommand = cli.Command{
	Name:      "build",
	ShortName: "b",
	Usage:     "generate a static website summarizing the compliance program",
	Action:    buildAction,
	Before:    beforeAll(pandocMustExist, cleanContainers),
}

func buildAction(c *cli.Context) error {
	err := render.Build("output", false)
	if err != nil {
		return errors.Wrap(err, "build failed")
	}
	return nil
}
