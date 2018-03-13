package model

import (
	"io/ioutil"
	"strings"

	"github.com/strongdm/comply/internal/path"
	"gopkg.in/yaml.v2"
)

// ReadPolicies loads policy records from the filesystem
func ReadPolicies() []Policy {
	var policies []Policy

	for _, f := range path.Policies() {
		p := Policy{}
		mdmd := loadMDMD(f.FullPath)
		yaml.Unmarshal([]byte(mdmd.yaml), &p)
		p.Body = mdmd.body
		policies = append(policies, p)
	}

	return policies
}

type metadataMarkdown struct {
	yaml string
	body string
}

func loadMDMD(path string) metadataMarkdown {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	content := string(bytes)
	components := strings.Split(content, "---")
	if len(components) == 1 {
		panic("Malformed metadata markdown file, must be of the form: YAML\\n---\\nmarkdown content")
	}
	yaml := components[0]
	body := strings.Join(components[1:], "---")
	return metadataMarkdown{yaml, body}
}
