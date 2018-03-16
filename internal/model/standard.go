package model

type Control struct {
	Family      string `yaml:"family"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type Standard struct {
	Name     string             `yaml:"name"`
	Controls map[string]Control `yaml:",inline"`
}
