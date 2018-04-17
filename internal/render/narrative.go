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
	"github.com/strongdm/comply/internal/config"
	"github.com/strongdm/comply/internal/model"
)

func renderNarrativeToDisk(wg *sync.WaitGroup, data *renderData, narrative *model.Narrative) {
	// only files that have been touched
	if !isNewer(narrative.FullPath, narrative.ModifiedAt) {
		return
	}
	recordModified(narrative.FullPath, narrative.ModifiedAt)

	ctx := context.Background()
	cli, err := client.NewEnvClient()
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

	wg.Add(1)
	go func(p *model.Narrative) {
		outputFilename := p.OutputFilename
		// save preprocessed markdown
		preprocessNarrative(data, p, filepath.Join(".", "output", outputFilename+".md"))

		cmd := []string{"--smart", "--toc", "-N", "--template=/source/templates/default.latex", "-o",
			fmt.Sprintf("/source/output/%s", outputFilename),
			fmt.Sprintf("/source/output/%s.md", outputFilename)}

		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: "strongdm/pandoc",
			Cmd:   cmd},
			hc, nil, "")
		if err != nil {
			panic(err)
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}

		cli.ContainerWait(ctx, resp.ID)

		_, err = cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
		// io.Copy(os.Stdout, rc)
		if err != nil {
			panic(err)
		}

		// remove preprocessed markdown
		os.Remove(filepath.Join(".", "output", outputFilename+".md"))
		wg.Done()
	}(narrative)
}

func preprocessNarrative(data *renderData, pol *model.Narrative, fullPath string) {
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
		panic(err)
	}
}
