package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var projectRoot string

var dockerAvailable, pandocAvailable bool

const (
	Jira      = "jira"
	GitHub    = "github"
	GitLab    = "gitlab"
	NoTickets = "none"
)

const (
	// UseDocker invokes pandoc within Docker
	UseDocker = "docker"
	// UsePandoc invokes pandoc directly
	UsePandoc = "pandoc"
)

// SetProjectRoot is used by the test suite.
func SetProjectRoot(dir string) {
	projectRoot = dir
}

type Project struct {
	Name           string                 `yaml:"name"`
	Pandoc         string                 `yaml:"pandoc,omitempty"`
	FilePrefix     string                 `yaml:"filePrefix"`
	Tickets        map[string]interface{} `yaml:"tickets"`
	ApprovedBranch string                 `yaml:"approvedBranch"`
}

// SetPandoc records pandoc availability during initialization
func SetPandoc(pandoc bool, docker bool) {
	pandocAvailable = pandoc
	dockerAvailable = docker
}

// WhichPandoc indicates which pandoc invocation path should be used
func WhichPandoc() string {
	cfg := Config()
	if cfg.Pandoc == UsePandoc {
		return UsePandoc
	}
	if cfg.Pandoc == UseDocker {
		return UseDocker
	}
	if pandocAvailable {
		return UsePandoc
	}
	return UseDocker
}

// YAML is the parsed contents of ProjectRoot()/config.yml.
func YAML() map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	cfgBytes, err := ioutil.ReadFile(filepath.Join(ProjectRoot(), "comply.yml"))
	if err != nil {
		panic("unable to load config.yml: " + err.Error())
	}
	yaml.Unmarshal(cfgBytes, &m)
	return m
}

// Exists tests for the presence of a comply configuration file.
func Exists() bool {
	_, err := ioutil.ReadFile(filepath.Join(ProjectRoot(), "comply.yml"))
	if err != nil {
		return false
	}
	return true
}

// Config is the parsed contents of ProjectRoot()/config.yml.
func Config() *Project {
	p := Project{}
	cfgBytes, err := ioutil.ReadFile(filepath.Join(ProjectRoot(), "comply.yml"))
	if err != nil {
		panic("unable to load config.yml: " + err.Error())
	}
	yaml.Unmarshal(cfgBytes, &p)
	return &p
}

// ProjectRoot is the fully-qualified path to the root directory.
func ProjectRoot() string {
	if projectRoot == "" {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		projectRoot = dir
	}

	return projectRoot
}

// TicketSystem indicates the type of the configured ticket system
func (p *Project) TicketSystem() (string, error) {
	if len(p.Tickets) > 1 {
		return NoTickets, errors.New("multiple ticket systems configured")
	}

	for k := range p.Tickets {
		switch k {
		case GitHub:
			return GitHub, nil
		case Jira:
			return Jira, nil
		case GitLab:
			return GitLab, nil
		case NoTickets:
			return NoTickets, nil
		default:
			// explicit error for this case
			return "", errors.New("unrecognized ticket system configured")
		}
	}

	// no ticket block configured
	return NoTickets, nil
}
