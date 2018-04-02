package cli

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/strongdm/comply/internal/config"
	"github.com/strongdm/comply/internal/plugin/github"
	"github.com/urfave/cli"
)

// Version is set by the build system
const Version = "0.0.0-development"

// Main should be invoked by the main function in the main package
func Main() {
	err := newApp().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "comply"
	app.Version = Version
	app.Usage = "policy compliance toolkit"

	app.Commands = []cli.Command{
		initCommand,
		buildCommand,
		todoCommand,
		serveCommand,
		syncCommand,
	}

	// Plugins
	github.Register()

	return app
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

func projectMustExist(c *cli.Context) error {
	_, err := ioutil.ReadFile(filepath.Join(config.ProjectRoot(), "comply.yml"))
	if err != nil {
		return errors.New("command must be run from the root of a valid comply project (comply.yml must exist; have you run `comply init`?)")
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

	// TODO: where does this const go?
	_, err = cli.ImagePull(ctx, "strongdm/pandoc", types.ImagePullOptions{})
	if err != nil {
		return dockerErr
	}
	return nil
}
