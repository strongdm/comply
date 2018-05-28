package render

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/config"
	"github.com/strongdm/comply/internal/model"
)

// TODO: refactor and eliminate duplication among narrative, policy renderers
func renderToFilesystem(wg *sync.WaitGroup, errOutputCh chan error, data *renderData, doc *model.Document, live bool) {
	// only files that have been touched
	if !isNewer(doc.FullPath, doc.ModifiedAt) {
		return
	}
	recordModified(doc.FullPath, doc.ModifiedAt)

	wg.Add(1)
	go func(p *model.Document) {
		defer wg.Done()

		outputFilename := p.OutputFilename
		// save preprocessed markdown
		err := preprocessDoc(data, p, filepath.Join(".", "output", outputFilename+".md"))
		if err != nil {
			errOutputCh <- errors.Wrap(err, "unable to preprocess")
			return
		}

		pandoc(outputFilename, errOutputCh)

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
	}(doc)
}

func preprocessDoc(data *renderData, pol *model.Document, fullPath string) error {
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
		return errors.Wrap(err, "unable to write preprocessed policy to disk")
	}
	return nil
}
