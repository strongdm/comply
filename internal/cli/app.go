package cli

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

// Version is set by the build system
const Version = "0.0.0-development"

// Main should be invoked by the main functino in the main package
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
		serveCommand,
	}

	return app
}
