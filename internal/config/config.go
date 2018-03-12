package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var projectRoot string

// SetProjectRoot is used by the test suite
func SetProjectRoot(dir string) {
	projectRoot = dir
}

// YAML is the parsed contents of ProjectRoot()/config.yml
func YAML() map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	cfgBytes, err := ioutil.ReadFile(filepath.Join(ProjectRoot(), "comply.yml"))
	if err != nil {
		panic("unable to load config.yml: " + err.Error())
	}
	yaml.Unmarshal(cfgBytes, &m)
	return m
}

// ProjectRoot is the fully-qualified path to the root directory
func ProjectRoot() string {
	if projectRoot == "" {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		projectRoot = dir
	}

	fullPath := filepath.Join(projectRoot, "comply.yml")
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("%s not found: comply must be run from the root directory of an initialized comply project", fullPath))
	}
	return projectRoot
}
