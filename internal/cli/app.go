package cli

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/config"
	"github.com/strongdm/comply/internal/gitlab"
	"github.com/strongdm/comply/internal/jira"
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
		beforeCommand(initCommand, notifyVersion),
	}

	app.Commands = append(app.Commands, beforeCommand(buildCommand, projectMustExist, notifyVersion))
	app.Commands = append(app.Commands, beforeCommand(procedureCommand, projectMustExist, notifyVersion))
	app.Commands = append(app.Commands, beforeCommand(schedulerCommand, projectMustExist, notifyVersion))
	app.Commands = append(app.Commands, beforeCommand(serveCommand, projectMustExist, notifyVersion))
	app.Commands = append(app.Commands, beforeCommand(syncCommand, projectMustExist, notifyVersion))
	app.Commands = append(app.Commands, beforeCommand(todoCommand, projectMustExist, notifyVersion))

	// Plugins
	github.Register()
	jira.Register()
	gitlab.Register()

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

// notifyVersion asynchronously notifies the availability of version updates
func notifyVersion(c *cli.Context) error {
	go func() {
		defer func() {
			recover() // suppress panic
		}()

		r, err := http.Get("http://comply-releases.s3.amazonaws.com/channel/stable/VERSION")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// fail silently
		}

		version := strings.TrimSpace(string(body))

		// only when numeric versions are present
		firstRune, _ := utf8.DecodeRuneInString(string(body))
		if unicode.IsDigit(firstRune) && version != Version {
			// only once every ~10 times
			if rand.Intn(10) == 0 {
				fmt.Fprintf(os.Stderr, "a new version of comply is available")
			}
		}
	}()
	return nil
}

func pandocMustExist(c *cli.Context) error {
	eitherMustExistErr := fmt.Errorf("\n\nPlease install either Docker or the pandoc package and re-run `%s`. Find OS-specific pandoc installation instructions at: [TODO]", c.Command.Name)

	pandocExistErr, found, goodVersion, pdfLatex := pandocBinaryMustExist(c)
	dockerExistErr, inPath, isRunning := dockerMustExist(c)

	config.SetPandoc(pandocExistErr == nil, dockerExistErr == nil)
	check := func(b bool) string {
		if b {
			return "✔"
		} else {
			return "✖"
		}

	}

	if pandocExistErr != nil && dockerExistErr != nil {

		fmt.Printf(`
[%s] pandoc binary installed and in PATH
[%s] pandoc version compatible
[%s] pdflatex binary installed and in PATH
[%s] docker binary installed
[%s] docker running

`, check(found), check(goodVersion), check(pdfLatex), check(inPath), check(isRunning))

		return eitherMustExistErr
	}

	// if we don't have pandoc, but we do have docker, execute a pull
	if (pandocExistErr != nil && dockerExistErr == nil) || config.WhichPandoc() == config.UseDocker {
		dockerPull(c)
	}

	return nil
}

func pandocBinaryMustExist(c *cli.Context) (e error, found, goodVersion, pdfLatex bool) {
	cmd := exec.Command("pandoc", "-v")
	outputRaw, err := cmd.Output()

	e = nil
	found = false
	goodVersion = false
	pdfLatex = false

	if err != nil {
		e = errors.Wrap(err, "error calling pandoc")
	} else {
		found = true
		goodVersion = true
		output := strings.TrimSpace((string(outputRaw)))
		versionErr := errors.New("cannot determine pandoc version")
		if !strings.HasPrefix(output, "pandoc") {
			e = versionErr
			goodVersion = false
		} else {
			re := regexp.MustCompile(`pandoc (\d+)\.(\d+)`)
			result := re.FindStringSubmatch(output)
			if len(result) != 3 {
				e = versionErr
				goodVersion = false
			} else {
				major, err := strconv.Atoi(result[1])
				if err != nil {
					e = versionErr
					goodVersion = false
				}
				minor, err := strconv.Atoi(result[2])
				if err != nil {
					e = versionErr
					goodVersion = false
				}
				if major < 2 || minor < 1 {
					e = errors.New("pandoc 2.1 or greater required")
					goodVersion = false
				}
			}
		}
	}

	// pdflatex must also be present
	cmd = exec.Command("pdflatex", "--version")
	outputRaw, err = cmd.Output()
	if err != nil {
		e = errors.Wrap(err, "error calling pdflatex")
	} else if !strings.Contains(string(outputRaw), "TeX") {
		e = errors.New("pdflatex is required")
	} else {
		pdfLatex = true
	}

	return e, found, goodVersion, pdfLatex
}

func dockerMustExist(c *cli.Context) (e error, inPath, isRunning bool) {
	dockerErr := fmt.Errorf("Docker must be available in order to run `%s`", c.Command.Name)

	inPath = true
	cmd := exec.Command("docker", "--version")
	_, err := cmd.Output()
	if err != nil {
		inPath = false
	}

	isRunning = true
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		isRunning = false
		return dockerErr, inPath, isRunning
	}

	_, err = cli.Ping(ctx)
	if err != nil {
		isRunning = false
		return dockerErr, inPath, isRunning
	}

	return nil, inPath, isRunning
}

func dockerPull(c *cli.Context) error {
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
