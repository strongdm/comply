package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
)

var initCommand = cli.Command{
	Name:   "init",
	Usage:  "initialize a new compliance repository (interactive)",
	Action: initAction,
}

func initAction(c *cli.Context) error {
	// create directory structure

	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("Invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Number",
		Validate: validate,
	}

	_, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil
	}

	chooser := promptui.Select{
		Label: "Compliance Regime",
		Items: []string{"SOC2", "Blank"},
	}

	_, _, err = chooser.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil
	}

	return nil
}
