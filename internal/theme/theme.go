package theme

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// SaveTo persists a compliance theme to a destination directory with optional
// template value replacements.
func SaveTo(themeName string, replace map[string]string, saveDir string) error {
	for _, name := range AssetNames() {
		prefix := themeName + "/"
		if strings.HasPrefix(name, prefix) {
			outputName := strings.TrimPrefix(name, prefix)
			assetDir, assetFilename := filepath.Split(outputName)
			err := os.MkdirAll(filepath.Join(saveDir, assetDir), os.FileMode(0755))
			if err != nil {
				return err
			}

			// special case for README.md and TODO.md: all other templates
			// are passed copied verbatim.
			if name == filepath.Join(themeName, "README.md") || name == filepath.Join(themeName, "TODO.md") {
				rootMdFile := string(MustAsset(name))

				var w bytes.Buffer
				var rootMdFileTemplate *template.Template
				rootMdFileTemplate, err = template.New("rootMdFile").Parse(rootMdFile)
				if err != nil {
					w.WriteString(fmt.Sprintf("# Error processing template:\n\n%s\n", err.Error()))
				} else {
					rootMdFileTemplate.Execute(&w, replace)
				}
				body := w.String()
				err = ioutil.WriteFile(filepath.Join(saveDir, assetDir, assetFilename), []byte(body), os.FileMode(0644))
			} else {
				err = ioutil.WriteFile(filepath.Join(saveDir, assetDir, assetFilename), MustAsset(name), os.FileMode(0644))
			}

			if err != nil {
				return err
			}
		}
	}
	// TODO
	return nil
}
