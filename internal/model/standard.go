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
func ControlsSatisfied(data *Data) map[string][]string {
	satisfied := make(map[string][]string)

	appendSatisfaction := func(in map[string][]string, k string, v string) []string {
		s, ok := in[k]
		if !ok {
			s = make([]string, 0)
		}
		s = append(s, v)
		return s
	}

	for _, n := range data.Narratives {
		for _, controlKeys := range n.Satisfies {
			for _, key := range controlKeys {
				satisfied[key] = appendSatisfaction(satisfied, key, n.OutputFilename)
			}
		}
	}
	for _, n := range data.Policies {
		for _, controlKeys := range n.Satisfies {
			for _, key := range controlKeys {
				satisfied[key] = appendSatisfaction(satisfied, key, n.OutputFilename)
			}
		}
	}
	for _, n := range data.Procedures {
		for _, controlKeys := range n.Satisfies {
			for _, key := range controlKeys {
				satisfied[key] = appendSatisfaction(satisfied, key, n.OutputFilename)
			}
		}
	}
	return satisfied
}
