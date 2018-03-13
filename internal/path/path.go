package path

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	FullPath string
	Info     os.FileInfo
}

// Policies lists all files in the policies directory
func Policies() []File {
	var filtered []File
	files, err := ioutil.ReadDir(filepath.Join(".", "policies"))
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		fmt.Println("CONSIDERING ", f)
		if !strings.HasSuffix(f.Name(), ".md") || strings.HasPrefix("README", strings.ToUpper(f.Name())) {
			continue
		}
		abs, err := filepath.Abs(filepath.Join(".", "policies", f.Name()))
		if err != nil {
			panic(err)
		}
		fmt.Println("appending")
		filtered = append(filtered, File{abs, f})
	}
	fmt.Println("returning from policies")
	return filtered
}
