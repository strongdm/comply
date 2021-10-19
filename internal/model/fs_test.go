package model

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/strongdm/comply/internal/path"
	"github.com/strongdm/comply/internal/util"
)

type ReadFiles struct{}

func beforeEach() {
	util.MockConfig()
}

func TestReadFiles(t *testing.T) {
	util.ExecuteTests(t, reflect.TypeOf(ReadFiles{}), beforeEach, nil)
}

// TestReadNarratives calls model.ReadNarratives checking for a valid return value.
func (tg ReadFiles) TestReadNarratives(t *testing.T) {
	filePath := fmt.Sprintf("%s/narratives/control.md", util.GetRootPath())
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
func (tg ReadFiles) TestReadNarrativesWhenThereAreNoNarratives(t *testing.T) {
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
func (tg ReadFiles) TestReadNarrativesFailsWhenInvalidNarrative(t *testing.T) {
	path.Narratives = func() ([]path.File, error) {
		filePath := fmt.Sprintf("%s/../fixtures/narratives/invalid-control.md", util.GetRootPath())
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
func (tg ReadFiles) TestReadProcedures(t *testing.T) {
	filePath := fmt.Sprintf("%s/procedures/workstation.md", util.GetRootPath())
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
func (tg ReadFiles) TestReadProceduresWhenThereAreNoProcedures(t *testing.T) {
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
func (tg ReadFiles) TestReadProceduresFailsWhenInvalidProcedure(t *testing.T) {
	path.Procedures = func() ([]path.File, error) {
		filePath := fmt.Sprintf("%s/../fixtures/procedures/invalid-workstation.md", util.GetRootPath())
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
func (tg ReadFiles) TestReadPolicies(t *testing.T) {
	filePath := fmt.Sprintf("%s/policies/access.md", util.GetRootPath())
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
func (tg ReadFiles) TestReadPoliciesWhenThereAreNoPolicies(t *testing.T) {
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
func (tg ReadFiles) TestReadPoliciesFailsWhenInvalidPolicy(t *testing.T) {
	path.Policies = func() ([]path.File, error) {
		filePath := fmt.Sprintf("%s/../fixtures/policies/invalid-access.md", util.GetRootPath())
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
func (tg ReadFiles) TestReadStandards(t *testing.T) {
	filePath := fmt.Sprintf("%s/standards/TSC-2017.yml", util.GetRootPath())
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
func (tg ReadFiles) TestReadStandardsWhenThereAreNoStandards(t *testing.T) {
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
func (tg ReadFiles) TestReadStandardsFailsWhenInvalidStandard(t *testing.T) {
	path.Standards = func() ([]path.File, error) {
		filePath := fmt.Sprintf("%s/../fixtures/standards/invalid-standard.yml", util.GetRootPath())
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
