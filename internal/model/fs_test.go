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
	filePath := fmt.Sprintf("%s/narratives/control.md", getRootPath())
	fileInfo, _ := os.Lstat(filePath)
	path.Narratives = func() ([]path.File, error) {
		return []path.File{
			{FullPath: filePath, Info: fileInfo},
		}, nil
	}

	documents, err := ReadNarratives()
	if err != nil {
		t.Fatalf(`ReadNarratives() returned an error %v`, err)
	}
	if len(documents) != 1 {
		t.Fatal(`Invalid number of documents`)
	}
	if documents[0].FullPath != filePath {
		t.Fatalf(`Invalid document path %s`, documents[0].FullPath)
	}
}

// TestReadNarrativesWhenThereAreNoNarratives calls model.ReadNarratives checking for a valid return when
// there are no narratives to process
func TestReadNarrativesWhenThereAreNoNarratives(t *testing.T) {
	mockConfig()
	path.Narratives = func() ([]path.File, error) {
		return []path.File{}, nil
	}

	documents, err := ReadNarratives()
	if err != nil {
		t.Fatalf(`ReadNarratives() returned an error %v`, err)
	}
	if len(documents) != 0 {
		t.Fatal(`Invalid number of documents`)
	}
}

// TestReadNarrativesFailsWhenInvalidNarrative calls model.ReadNarratives checking for an error return when
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

// TestReadProcedures calls model.ReadProcedures checking for a valid return value.
func TestReadProcedures(t *testing.T) {
	mockConfig()
	filePath := fmt.Sprintf("%s/procedures/workstation.md", getRootPath())
	fileInfo, _ := os.Lstat(filePath)
	path.Procedures = func() ([]path.File, error) {
		return []path.File{
			{FullPath: filePath, Info: fileInfo},
		}, nil
	}

	documents, err := ReadProcedures()
	if err != nil {
		t.Fatalf(`ReadProcedures() returned an error %v`, err)
	}
	if len(documents) != 1 {
		t.Fatal(`Invalid number of documents`)
	}
	if documents[0].FullPath != filePath {
		t.Fatalf(`Invalid document path %s`, documents[0].FullPath)
	}
}

// TestReadProceduresWhenThereAreNoProcedures calls model.ReadProcedures checking for a valid return when
// there are no procedures to process
func TestReadProceduresWhenThereAreNoProcedures(t *testing.T) {
	mockConfig()
	path.Procedures = func() ([]path.File, error) {
		return []path.File{}, nil
	}

	documents, err := ReadProcedures()
	if err != nil {
		t.Fatalf(`ReadProcedures() returned an error %v`, err)
	}
	if len(documents) != 0 {
		t.Fatal(`Invalid number of documents`)
	}
}

// TestReadProceduresFailsWhenInvalidProcedure calls model.ReadProcedures checking for an error return when
// there is an invalid procedure
func TestReadProceduresFailsWhenInvalidProcedure(t *testing.T) {
	mockConfig()
	path.Procedures = func() ([]path.File, error) {
		filePath := fmt.Sprintf("%s/../fixtures/procedures/invalid-workstation.md", getRootPath())
		fileInfo, _ := os.Lstat(filePath)
		return []path.File{
			{FullPath: filePath, Info: fileInfo},
		}, nil
	}

	_, err := ReadProcedures()
	if err == nil {
		t.Fatal(`ReadProcedures() was expected to fail`, err)
	}
}

// TestReadPolicies calls model.ReadPolicies checking for a valid return value.
func TestReadPolicies(t *testing.T) {
	mockConfig()
	filePath := fmt.Sprintf("%s/policies/access.md", getRootPath())
	fileInfo, _ := os.Lstat(filePath)
	path.Policies = func() ([]path.File, error) {
		return []path.File{
			{FullPath: filePath, Info: fileInfo},
		}, nil
	}

	documents, err := ReadPolicies()
	if err != nil {
		t.Fatalf(`ReadPolicies() returned an error %v`, err)
	}
	if len(documents) != 1 {
		t.Fatal(`Invalid number of documents`)
	}
	if documents[0].FullPath != filePath {
		t.Fatalf(`Invalid document path %s`, documents[0].FullPath)
	}
}

// TestReadPoliciesWhenThereAreNoPolicies calls model.ReadPolicies checking for a valid return when
// there are no policies to process
func TestReadPoliciesWhenThereAreNoPolicies(t *testing.T) {
	mockConfig()
	path.Policies = func() ([]path.File, error) {
		return []path.File{}, nil
	}

	documents, err := ReadPolicies()
	if err != nil {
		t.Fatalf(`ReadPolicies() returned an error %v`, err)
	}
	if len(documents) != 0 {
		t.Fatal(`Invalid number of documents`)
	}
}

// TestReadPoliciesFailsWhenInvalidPolicy calls model.ReadPolicies checking for an error return when
// there is an invalid policy
func TestReadPoliciesFailsWhenInvalidPolicy(t *testing.T) {
	mockConfig()
	path.Policies = func() ([]path.File, error) {
		filePath := fmt.Sprintf("%s/../fixtures/policies/invalid-access.md", getRootPath())
		fileInfo, _ := os.Lstat(filePath)
		return []path.File{
			{FullPath: filePath, Info: fileInfo},
		}, nil
	}

	_, err := ReadPolicies()
	if err == nil {
		t.Fatal(`ReadPolicies() was expected to fail`, err)
	}
}

// TestReadStandards calls model.ReadStandards checking for a valid return value.
func TestReadStandards(t *testing.T) {
	mockConfig()
	filePath := fmt.Sprintf("%s/standards/TSC-2017.yml", getRootPath())
	fileInfo, _ := os.Lstat(filePath)
	path.Standards = func() ([]path.File, error) {
		return []path.File{
			{FullPath: filePath, Info: fileInfo},
		}, nil
	}

	documents, err := ReadStandards()
	if err != nil {
		t.Fatalf(`ReadStandards() returned an error %v`, err)
	}
	if len(documents) != 1 {
		t.Fatal(`Invalid number of documents`)
	}
}

// TestReadStandardsWhenThereAreNoStandards calls model.ReadStandards checking for a valid return when
// there are no standards to process
func TestReadStandardsWhenThereAreNoStandards(t *testing.T) {
	mockConfig()
	path.Standards = func() ([]path.File, error) {
		return []path.File{}, nil
	}

	documents, err := ReadStandards()
	if err != nil {
		t.Fatalf(`ReadStandards() returned an error %v`, err)
	}
	if len(documents) != 0 {
		t.Fatal(`Invalid number of documents`)
	}
}

// TestReadStandardsFailsWhenInvalidStandard calls model.ReadStandards checking for an error return when
// there is an invalid standard
func TestReadStandardsFailsWhenInvalidStandard(t *testing.T) {
	mockConfig()
	path.Standards = func() ([]path.File, error) {
		filePath := fmt.Sprintf("%s/../fixtures/standards/invalid-standard.yml", getRootPath())
		fileInfo, _ := os.Lstat(filePath)
		return []path.File{
			{FullPath: filePath, Info: fileInfo},
		}, nil
	}

	_, err := ReadStandards()
	if err == nil {
		t.Fatal(`ReadStandards() was expected to fail`, err)
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
