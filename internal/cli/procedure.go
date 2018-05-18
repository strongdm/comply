package cli

import (
	"fmt"

	"github.com/strongdm/comply/internal/model"
	"github.com/urfave/cli"
)

var procedureCommand = cli.Command{
	Name:      "procedure",
	ShortName: "proc",
	Usage:     "create ticket by procedure ID",
	ArgsUsage: "procedureID",
	Action:    procedureAction,
	Before:    projectMustExist,
}

func procedureAction(c *cli.Context) error {
	procedures, err := model.ReadProcedures()
	if err != nil {
		return err
	}

	if c.NArg() != 1 {
		return cli.NewExitError("provide a procedure ID", 1)
	}

	procedureID := c.Args().First()

	for _, procedure := range procedures {
		if procedure.ID == procedureID {
			// TODO: don't hardcode GH
			tp := model.GetPlugin(model.GitHub)
			tp.Create(&model.Ticket{
				Name: procedure.Name,
				Body: fmt.Sprintf("%s\n\n\n---\nProcedure-ID: %s", procedure.Body, procedure.ID),
			}, []string{"comply", "comply-procedure"})
			return nil
		}
	}

	return cli.NewExitError(fmt.Sprintf("unknown procedure ID: %s", procedureID), 1)
}
