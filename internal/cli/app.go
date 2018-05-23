package cli

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/config"
	"github.com/strongdm/comply/internal/plugin/github"
	"github.com/urfave/cli"
)

// Version is set by the build system.
var Version = ""

// Main should be invoked by the main function in the main package.
func Main() {
	err := newApp().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "comply"
	if Version == "" {
		app.HideVersion = true
	}
	app.Version = Version
	app.Usage = "policy compliance toolkit"

	app.Commands = []cli.Command{
		initCommand,
	}

	app.Commands = append(app.Commands, beforeCommand(buildCommand, projectMustExist))
	app.Commands = append(app.Commands, beforeCommand(procedureCommand, projectMustExist))
	app.Commands = append(app.Commands, beforeCommand(schedulerCommand, projectMustExist))
	app.Commands = append(app.Commands, beforeCommand(serveCommand, projectMustExist))
	app.Commands = append(app.Commands, beforeCommand(syncCommand, projectMustExist))
	app.Commands = append(app.Commands, beforeCommand(todoCommand, projectMustExist))

	// Plugins
	github.Register()

	return app
}

func beforeCommand(c cli.Command, bf ...cli.BeforeFunc) cli.Command {
	if c.Before == nil {
		c.Before = beforeAll(bf...)
	} else {
		c.Before = beforeAll(append(bf, c.Before)...)
	}
	return c
}

func beforeAll(bf ...cli.BeforeFunc) cli.BeforeFunc {
	return func(c *cli.Context) error {
		for _, f := range bf {
			if err := f(c); err != nil {
				return err
			}
		}
		return nil
	}
}

func feedbackError(message string) error {
	return errors.New(fmt.Sprintf("\n\nERROR\n=====\n%s\n", message))
}

func projectMustExist(c *cli.Context) error {
	_, err := ioutil.ReadFile(filepath.Join(config.ProjectRoot(), "comply.yml"))
	if err != nil {
		return feedbackError("command must be run from the root of a valid comply project (comply.yml must exist; have you run `comply init`?)")
	}
	return nil
}

func ticketingMustBeConfigured(c *cli.Context) error {
	p := config.Config()
	if p.Tickets == nil || len(p.Tickets) != 1 {
		return feedbackError("comply.yml must contain a valid ticketing configuration")
	}
	return nil
}

func pandocMustExist(c *cli.Context) error {
	eitherMustExistErr := fmt.Errorf("Please install either Docker or the pandoc package and re-run `%s`", c.Command.Name)

	pandocExistErr := pandocBinaryMustExist(c)
	dockerExistErr := dockerMustExist(c)
	config.SetPandoc(pandocExistErr == nil, dockerExistErr == nil)

	if pandocExistErr != nil && dockerExistErr != nil {
		return eitherMustExistErr
	}

	return nil
}

func pandocBinaryMustExist(c *cli.Context) error {
	cmd := exec.Command("pandoc", "-v")
	outputRaw, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "error calling pandoc")
	}

	output := strings.TrimSpace((string(outputRaw)))
	versionErr := errors.New("cannot determine pandoc version")
	if !strings.HasPrefix(output, "pandoc") {
		return versionErr
	}

	re := regexp.MustCompile(`pandoc (\d+)\.(\d+)`)
	result := re.FindStringSubmatch(output)
	if len(result) != 3 {
		return versionErr
	}

	major, err := strconv.Atoi(result[1])
	if err != nil {
		return versionErr
	}
	minor, err := strconv.Atoi(result[2])
	if err != nil {
		return versionErr
	}

	if major < 2 || minor < 1 {
		return errors.New("pandoc 2.1 or greater required")
	}
	return nil
}

func dockerMustExist(c *cli.Context) error {
	dockerErr := fmt.Errorf("Docker must be available in order to run `%s`", c.Command.Name)

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return dockerErr
	}

	done := make(chan struct{})
	defer close(done)

	go func() {
		// if docker IO takes more than N seconds, notify user we're (likely) downloading the pandoc image
		longishPull := time.After(time.Second * 5)

		select {
		case <-longishPull:
			fmt.Print("Pulling strongdm/pandoc:latest Docker image (this will take some time) ")

			go func() {
				for {
					fmt.Print(".")
					select {
					case <-done:
						fmt.Print(" done.\n")
						return
					default:
						time.Sleep(1 * time.Second)
					}
				}
			}()
		case <-done:
			// in this case, the docker pull was quick -- suggesting we already have the container
		}
	}()

	r, err := cli.ImagePull(ctx, "strongdm/pandoc:latest", types.ImagePullOptions{})
	if err != nil {
		return dockerErr
	}
	defer r.Close()

	// hold function open until all docker IO is complete
	io.Copy(ioutil.Discard, r)

	return nil
}

func cleanContainers(c *cli.Context) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		// no Docker? nothing to clean.
		return nil
	}

	_, err = cli.Ping(ctx)
	if err != nil {
		// no Docker? nothing to clean.
		return nil
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return errors.Wrap(err, "error listing containers during cleanup")
	}

	for _, c := range containers {
		// assume this container was leftover from previous aborted run
		if strings.HasPrefix(c.Image, "strongdm/pandoc") {
			d := time.Second * 2
			err = cli.ContainerStop(ctx, c.ID, &d)
			if err != nil {
				fmt.Printf("Unable to stop container ID %s\n", c.ID)
			}

			err = cli.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{Force: true})
			if err != nil {
				fmt.Printf("Unable to remove container ID %s, please attempt manual removal\n", c.ID)
			}
		}
	}
	return nil
}
