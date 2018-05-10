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
	"github.com/strongdm/comply/internal/model"
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
		Label:    "Organization Name",
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
		Label: "Compliance Templates",
		Items: []string{"SOC2", "Blank"},
	}

	choice, _, err := chooser.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	themeName := "comply-blank"
	switch choice {
	case 0:
		themeName = "comply-soc2"
	case 1:
		themeName = "comply-blank"
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

	ticketing := model.Github

	switch choice {
	case 0:
		ticketing = model.Github
	case 1:
		ticketing = model.JIRA
	default:
		panic("unrecognized selection")
	}

	ticketConfig := make(map[string]string)

	plugin := model.GetPlugin(ticketing)
	ticketPrompts := plugin.Prompts()
	for k, prompt := range ticketPrompts {
		p := promptui.Prompt{
			Label:    prompt,
			Validate: atLeast(2),
		}

		v, err := p.Run()
		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			return err
		}
		ticketConfig[k] = v
	}

	p := config.Project{}
	p.Name = name
	p.FilePrefix = prefix
	p.Tickets = make(map[string]interface{})
	p.Tickets[string(ticketing)] = ticketConfig

	x, _ := yaml.Marshal(&p)
	err = ioutil.WriteFile(filepath.Join(config.ProjectRoot(), "comply.yml"), x, os.FileMode(0644))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = theme.SaveTo(themeName, config.ProjectRoot())
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	success := fmt.Sprintf("%s Compliance initialized successfully", name)
	fmt.Println(strings.Repeat("=", len(success)+2))
	fmt.Printf("%s %s\n", promptui.IconGood, success)

	return nil
}
