package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/strongdm/comply/internal/config"
	"github.com/strongdm/comply/internal/path"
	"gopkg.in/yaml.v2"
)

// TestReadNarratives calls model.ReadNarratives checking for a valid return value.
func TestReadNarratives(t *testing.T) {
	mockConfig()
	path.Narratives = func() ([]path.File, error) {
		filePath := fmt.Sprintf("%s/narratives/control.md", getRootPath())
		fileInfo, _ := os.Lstat(filePath)
		return []path.File{
			{FullPath: filePath, Info: fileInfo},
		}, nil
	}

	_, err := ReadNarratives()
	if err != nil {
		t.Fatalf(`ReadNarratives() returned an error %v`, err)
	}
}

// TestReadNarratives calls model.ReadNarratives checking for a valid return when
// there are no narratives to process
func TestReadNarrativesWhenThereAreNoNarratives(t *testing.T) {
	mockConfig()
	path.Narratives = func() ([]path.File, error) {
		return []path.File{}, nil
	}

	_, err := ReadNarratives()
	if err != nil {
		t.Fatalf(`ReadNarratives() returned an error %v`, err)
	}
}

// TestReadNarratives calls model.ReadNarratives checking for an error return when
// there is an invalid narrative
func TestReadNarrativesFailsWhenInvalidNarrative(t *testing.T) {
	mockConfig()
	path.Narratives = func() ([]path.File, error) {
		filePath := fmt.Sprintf("%s/../fixtures/narratives/invalid-control.md", getRootPath())
		fileInfo, _ := os.Lstat(filePath)
		return []path.File{
			{FullPath: filePath, Info: fileInfo},
		}, nil
	}

	_, err := ReadNarratives()
	if err == nil {
		t.Fatal(`ReadNarratives() was expected to fail`)
	}
}

func mockConfig() {
	config.Config = func() *config.Project {
		p := config.Project{}
		cfgBytes, _ := ioutil.ReadFile(filepath.Join(getRootPath(), "comply.yml.example"))
		yaml.Unmarshal(cfgBytes, &p)
		return &p
	}
}

func getRootPath() string {
	_, fileName, _, _ := runtime.Caller(0)
	fileDir := filepath.Dir(fileName)
	return fmt.Sprintf("%s/../../example", fileDir)
}
