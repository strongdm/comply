package render

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/config"
	"github.com/strongdm/comply/internal/model"
)

// TODO: refactor and eliminate duplication among narrative, policy renderers
func renderNarrativeToDisk(wg *sync.WaitGroup, errOutputCh chan error, data *renderData, narrative *model.Narrative, live bool) {
	// only files that have been touched
	if !isNewer(narrative.FullPath, narrative.ModifiedAt) {
		return
	}
	recordModified(narrative.FullPath, narrative.ModifiedAt)

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

	wg.Add(1)
	go func(p *model.Narrative) {
		defer wg.Done()

		outputFilename := p.OutputFilename
		// save preprocessed markdown
		err = preprocessNarrative(data, p, filepath.Join(".", "output", outputFilename+".md"))
		if err != nil {
			errOutputCh <- errors.Wrap(err, "unable to preprocess")
			return
		}

		cmd := []string{"--smart", "--toc", "-N", "--template=/source/templates/default.latex", "-o",
			fmt.Sprintf("/source/output/%s", outputFilename),
			fmt.Sprintf("/source/output/%s.md", outputFilename)}

		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: "strongdm/pandoc",
			Cmd:   cmd},
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

		// remove preprocessed markdown
		err = os.Remove(filepath.Join(".", "output", outputFilename+".md"))
		if err != nil {
			errOutputCh <- err
			return
		}

		rel, err := filepath.Rel(config.ProjectRoot(), p.FullPath)
		if err != nil {
			rel = p.FullPath
		}
		fmt.Printf("%s -> %s\n", rel, filepath.Join("output", p.OutputFilename))

	}(narrative)
}

func preprocessNarrative(data *renderData, pol *model.Narrative, fullPath string) error {
	cfg := config.Config()

	var w bytes.Buffer
	bodyTemplate, err := template.New("body").Parse(pol.Body)
	if err != nil {
		w.WriteString(fmt.Sprintf("# Error processing template:\n\n%s\n", err.Error()))
	} else {
		bodyTemplate.Execute(&w, data)
	}
	body := w.String()

	revisionTable := ""
	satisfiesTable := ""

	// ||Date|Comment|
	// |---+------|
	// | 4 Jan 2018 | Initial Version |
	// Table: Document history

	if len(pol.Satisfies) > 0 {
		rows := ""
		for standard, keys := range pol.Satisfies {
			rows += fmt.Sprintf("| %s | %s |\n", standard, strings.Join(keys, ", "))
		}
		satisfiesTable = fmt.Sprintf("|Standard|Controls Satisfied|\n|-------+--------------------------------------------|\n%s\nTable: Control satisfaction\n", rows)
	}

	if len(pol.Revisions) > 0 {
		rows := ""
		for _, rev := range pol.Revisions {
			rows += fmt.Sprintf("| %s | %s |\n", rev.Date, rev.Comment)
		}
		revisionTable = fmt.Sprintf("|Date|Comment|\n|---+--------------------------------------------|\n%s\nTable: Document history\n", rows)
	}

	doc := fmt.Sprintf(`%% %s
%% %s
%% %s

---
header-includes: yes
head-content: "%s"
foot-content: "%s confidential %d"
---

%s

%s

\newpage
%s`,
		pol.Name,
		cfg.Name,
		fmt.Sprintf("%s %d", pol.ModifiedAt.Month().String(), pol.ModifiedAt.Year()),
		pol.Name,
		cfg.Name,
		time.Now().Year(),
		satisfiesTable,
		revisionTable,
		body,
	)
	err = ioutil.WriteFile(fullPath, []byte(doc), os.FileMode(0644))
	if err != nil {
		return errors.Wrap(err, "unable to write preprocessed narrative to disk")
	}
	return nil
}
