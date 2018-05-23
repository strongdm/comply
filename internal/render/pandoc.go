package render

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/config"
)

func pandoc(outputFilename string, errOutputCh chan error) {
	if config.WhichPandoc() == config.UsePandoc {
		err := pandocPandoc(outputFilename)
		if err != nil {
			errOutputCh <- err
		}
	} else {
		dockerPandoc(outputFilename, errOutputCh)
	}
}

func dockerPandoc(outputFilename string, errOutputCh chan error) {
	// TODO: switch to new args once docker image is updated
	// cmd21 := []string{"-f", "markdown+smart", "--toc", "-N", "--template", "templates/default.latex", "-o", fmt.Sprintf("output/%s", outputFilename), fmt.Sprintf("output/%s.md", outputFilename)}
	cmd19 := []string{"--smart", "--toc", "-N", "--template=/source/templates/default.latex", "-o",
		fmt.Sprintf("/source/output/%s", outputFilename),
		fmt.Sprintf("/source/output/%s.md", outputFilename)}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		errOutputCh <- errors.Wrap(err, "unable to read Docker environment")
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		errOutputCh <- errors.Wrap(err, "unable to get workding directory")
		return
	}

	hc := &container.HostConfig{
		Binds: []string{pwd + ":/source"},
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "strongdm/pandoc",
		Cmd:   cmd19},
		hc, nil, "")

	if err != nil {
		errOutputCh <- errors.Wrap(err, "unable to create Docker container")
		return
	}

	defer func() {
		timeout := 2 * time.Second
		cli.ContainerStop(ctx, resp.ID, &timeout)
		err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			errOutputCh <- errors.Wrap(err, "unable to remove container")
			return
		}
	}()

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		errOutputCh <- errors.Wrap(err, "unable to start Docker container")
		return
	}

	_, err = cli.ContainerWait(ctx, resp.ID)
	if err != nil {
		errOutputCh <- errors.Wrap(err, "error awaiting Docker container")
		return
	}

	_, err = cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		errOutputCh <- errors.Wrap(err, "error reading Docker container logs")
		return
	}
}

// 🐼
func pandocPandoc(outputFilename string) error {
	// -f markdown+smart --toc -N --template=templates/default.latex -o output/%s output/%s.md
	cmd := exec.Command("pandoc", "-f", "markdown+smart", "--toc", "-N", "--template", "templates/default.latex", "-o", fmt.Sprintf("output/%s", outputFilename), fmt.Sprintf("output/%s.md", outputFilename))
	outputRaw, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(outputRaw))
		return errors.Wrap(err, "error calling pandoc")
	}
	return nil
}
