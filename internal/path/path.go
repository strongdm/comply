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

// Standards lists all standard files
func Standards() []File {
	return filesFor("standards", "yml")
}

// Narratives lists all policy files
func Narratives() []File {
	return filesFor("narratives", "md")
}

// Policies lists all policy files
func Policies() []File {
	return filesFor("policies", "md")
}

// Procedures lists all procedure files
func Procedures() []File {
	return filesFor("procedures", "md")
}

func filesFor(name, extension string) []File {
	var filtered []File
	files, err := ioutil.ReadDir(filepath.Join(".", name))
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), "."+extension) || strings.HasPrefix(strings.ToUpper(f.Name()), "README") {
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
