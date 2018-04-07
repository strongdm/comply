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

// ControlsSatisfied determines the unique controls currently satisfied by all Narratives, Policies, and Procedures
func ControlsSatisfied(data *Data) map[string]bool {
	satisfied := make(map[string]bool)
	for _, n := range data.Narratives {
		for _, controlKeys := range n.Satisfies {
			for _, key := range controlKeys {
				satisfied[key] = true
			}
		}
	}
	for _, n := range data.Policies {
		for _, controlKeys := range n.Satisfies {
			for _, key := range controlKeys {
				satisfied[key] = true
			}
		}
	}
	for _, n := range data.Procedures {
		for _, controlKeys := range n.Satisfies {
			for _, key := range controlKeys {
				satisfied[key] = true
			}
		}
	}
	return satisfied
}
