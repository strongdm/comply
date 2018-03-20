package theme

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func SaveTo(themeName, saveDir string) error {
	for _, name := range AssetNames() {
		prefix := themeName + "/"
		if strings.HasPrefix(name, prefix) {
			outputName := strings.TrimPrefix(name, prefix)
			assetDir, assetFilename := filepath.Split(outputName)
			err := os.MkdirAll(filepath.Join(saveDir, assetDir), os.FileMode(0755))
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(filepath.Join(saveDir, assetDir, assetFilename), MustAsset(name), os.FileMode(0644))
			if err != nil {
				return err
			}
		}
	}
	// TODO
	return nil
}
