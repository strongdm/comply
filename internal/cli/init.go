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

const whatNow = `Next steps:

	* Customize this directory using a text editor ('cat TODO.md' for ideas)
	* Try 'comply build' and 'comply serve'
	* View output/index.html
	* Add this directory to source control
`

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
		Label:    "PDF Filename Prefix (no spaces, no trailing separator)",
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

	fmt.Printf("\nComply relies on your ticketing system for optional procedure tracking. You can always come back and enable this integration later.\n\n\n")

	chooser = promptui.Select{
		Label: "Ticket System",
		Items: []string{"GitHub", "Jira", "GitLab", "None"},
	}

	choice, _, err = chooser.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	var tickets map[string]interface{}
	ticketing := model.GitHub

	switch choice {
	case 0:
		ticketing = model.GitHub
	case 1:
		ticketing = model.Jira
	case 2:
		ticketing = model.GitLab
	default:
		ticketing = model.NoTickets
	}

	if ticketing != model.NoTickets {
		chooser = promptui.Select{
			Label: "Configure ticketing system?",
			Items: []string{fmt.Sprintf("Configure %s now", string(ticketing)), "Configure later (via comply.yml)"},
		}

		choice, _, err = chooser.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return err
		}

		tickets = make(map[string]interface{})
		ticketConfig := make(map[string]string)
		tickets[string(ticketing)] = ticketConfig

		if choice == 0 {
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
		}
	}

	p := config.Project{}
	p.Name = name
	p.FilePrefix = prefix
	p.Tickets = tickets

	x, _ := yaml.Marshal(&p)
	err = ioutil.WriteFile(filepath.Join(config.ProjectRoot(), "comply.yml"), x, os.FileMode(0644))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	replace := make(map[string]string)
	replace["Name"] = p.Name

	err = theme.SaveTo(themeName, replace, config.ProjectRoot())
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	success := fmt.Sprintf("%s Compliance initialized successfully!", name)
	fmt.Printf("%s %s\n\n", promptui.IconGood, success)
	fmt.Println(whatNow)

	return nil
}
