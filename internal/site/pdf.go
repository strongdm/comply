package site

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/strongdm/comply/internal/config"
	"github.com/strongdm/comply/internal/model"
)

func pdf(output string, live bool, wg *sync.WaitGroup) {
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

	for {
		var pwg sync.WaitGroup
		for _, policy := range model.ReadPolicies() {
			// only files that have been touched
			if !isNewer(policy.FullPath, policy.ModifiedAt) {
				continue
			}
			recordModified(policy.FullPath, policy.ModifiedAt)

			pwg.Add(1)
			go func(p *model.Policy) {
				outputFilename := p.OutputFilename
				// save preprocessed markdown
				preprocessPandoc(p, filepath.Join(".", "output", outputFilename+".md"))

				resp, err := cli.ContainerCreate(ctx, &container.Config{
					Image: "jagregory/pandoc",
					Cmd: []string{"--smart", "--toc", "-N", "--template=/source/templates/default.latex", "-o",
						fmt.Sprintf("/source/output/%s", outputFilename),
						fmt.Sprintf("/source/output/%s.md", outputFilename),
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

				// remove preprocessed markdown
				os.Remove(filepath.Join(".", "output", outputFilename+".md"))
				pwg.Done()
			}(policy)
		}

		pwg.Wait()

		if !live {
			wg.Done()
			return
		}
		<-subscribe()
	}
}

func preprocessPandoc(pol *model.Policy, fullPath string) {
	cfg := config.Config()
	doc := fmt.Sprintf("%% %s\n%% %s\n%% %s\n\n%s",
		pol.Name,
		cfg.Name,
		fmt.Sprintf("%s %d", pol.ModifiedAt.Month().String(), pol.ModifiedAt.Year()),
		pol.Body)
	err := ioutil.WriteFile(fullPath, []byte(doc), os.FileMode(0644))
	if err != nil {
		panic(err)
	}
}
