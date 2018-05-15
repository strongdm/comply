package path

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// File wraps an os.FileInfo as well as the absolute path to the underlying file.
type File struct {
	FullPath string
	Info     os.FileInfo
}

// Standards lists all standard files.
func Standards() ([]File, error) {
	return filesFor("standards", "yml")
}

// Narratives lists all narrative files.
func Narratives() ([]File, error) {
	return filesFor("narratives", "md")
}

// Policies lists all policy files.
func Policies() ([]File, error) {
	return filesFor("policies", "md")
}

// Procedures lists all procedure files.
func Procedures() ([]File, error) {
	return filesFor("procedures", "md")
}

func filesFor(name, extension string) ([]File, error) {
	var filtered []File
	files, err := ioutil.ReadDir(filepath.Join(".", name))
	if err != nil {
		return nil, errors.Wrap(err, "unable to load files for: "+name)
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), "."+extension) || strings.HasPrefix(strings.ToUpper(f.Name()), "README") {
			continue
		}
		abs, err := filepath.Abs(filepath.Join(".", name, f.Name()))
		if err != nil {
			return nil, errors.Wrap(err, "unable to load file: "+f.Name())
		}
		filtered = append(filtered, File{abs, f})
	}
	return filtered, nil
}
