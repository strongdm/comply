package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/strongdm/comply/internal/config"
	"github.com/strongdm/comply/internal/theme"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

var initCommand = cli.Command{
	Name:   "init",
	Usage:  "initialize a new compliance repository (interactive)",
	Action: initAction,
}

func initAction(c *cli.Context) error {
	fi, _ := ioutil.ReadDir(config.ProjectRoot())
	if len(fi) > 0 {
		return errors.New("init must be run from an empty directory")
	}

	atLeast := func(n int) func(string) error {
		return func(input string) error {
			if len(input) < n {
				return errors.New("Too short")
			}
			return nil
		}
	}

	noSpaces := func(s string) error {
		if strings.ContainsAny(s, "\n\t ") {
			return errors.New("Must not contain spaces")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Project Name",
		Validate: atLeast(1),
	}

	name, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	prompt = promptui.Prompt{
		Label:    "PDF Filename Prefix",
		Default:  strings.Split(name, " ")[0],
		Validate: noSpaces,
	}

	prefix, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	chooser := promptui.Select{
		Label: "Project Theme",
		Items: []string{"SOC2", "Blank"},
	}

	choice, _, err := chooser.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	themeName := "blank"
	switch choice {
	case 0:
		themeName = "soc2"
	case 1:
		themeName = "blank"
	default:
		panic("unrecognized selection")
	}

	chooser = promptui.Select{
		Label: "Ticket System",
		Items: []string{"Github", "JIRA"},
	}

	choice, _, err = chooser.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	ticketing := "github"
	switch choice {
	case 0:
		ticketing = "github"
	case 1:
		ticketing = "jira"
	default:
		panic("unrecognized selection")
	}

	p := config.Project{}
	p.Name = name
	p.FilePrefix = prefix
	p.Tickets = make(map[string]interface{})
	p.Tickets[ticketing] = "see documentation for format"

	x, _ := yaml.Marshal(&p)
	ioutil.WriteFile(filepath.Join(config.ProjectRoot(), "comply.yml"), x, os.FileMode(0644))

	theme.SaveTo(themeName, config.ProjectRoot())
	return nil
}
