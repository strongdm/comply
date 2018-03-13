package path

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	FullPath string
	Info     os.FileInfo
}

// Policies lists all policy files
func Policies() []File {
	return filesForDir("policies")
}

// Procedures lists all procedure files
func Procedures() []File {
	return filesForDir("procedures")
}

func filesForDir(name string) []File {
	var filtered []File
	files, err := ioutil.ReadDir(filepath.Join(".", name))
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".md") || strings.HasPrefix(strings.ToUpper(f.Name()), "README") {
			continue
		}
		abs, err := filepath.Abs(filepath.Join(".", name, f.Name()))
		if err != nil {
			panic(err)
		}
		filtered = append(filtered, File{abs, f})
	}
	return filtered
}
