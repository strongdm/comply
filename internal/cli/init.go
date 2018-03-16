package cli

import (
	"fmt"

	"github.com/urfave/cli"
)

var initCommand = cli.Command{
	Name:   "init",
	Usage:  "initialize a new compliance repository (interactive)",
	Action: initAction,
}

func initAction(c *cli.Context) error {
	fmt.Println("boom, you've initialized a new compliance repository")
	// create directory structure
	return nil
}
