package site

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func pdf(output string, live bool, wg *sync.WaitGroup) {
	for {
		files, err := ioutil.ReadDir("./policies")
		if err != nil {
			panic(err)
		}

		for _, fileInfo := range files {
			// only non-README markdown files
			if !strings.HasSuffix(fileInfo.Name(), ".md") || strings.HasPrefix("README", strings.ToUpper(fileInfo.Name())) {
				continue
			}

			// only files that have been touched
			if !isNewer(fileInfo) {
				continue
			}
			recordModified(fileInfo)

			basename := strings.Replace(fileInfo.Name(), ".md", "", -1)

			ctx := context.Background()
			cli, err := client.NewEnvClient()
			if err != nil {
				panic(err)
			}

			_, err = cli.ImagePull(ctx, "jagregory/pandoc", types.ImagePullOptions{})
			if err != nil {
				panic(err)
			}

			pwd, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			hc := &container.HostConfig{
				Binds: []string{pwd + ":/source"},
			}

			resp, err := cli.ContainerCreate(ctx, &container.Config{
				Image: "jagregory/pandoc",
				Cmd: []string{"--smart", "--toc", "-N", "--template=/source/templates/default.latex", "-o",
					fmt.Sprintf("/source/output/%s.pdf", basename),
					fmt.Sprintf("/source/policies/%s.md", basename),
				},
			}, hc, nil, "")
			if err != nil {
				panic(err)
			}

			if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
				panic(err)
			}

			cli.ContainerWait(ctx, resp.ID)

			_, err = cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
			if err != nil {
				panic(err)
			}

			// io.Copy(os.Stdout, out)
		}

		if !live {
			wg.Done()
			return
		}
		<-subscribe()
	}
}
