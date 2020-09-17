package model

type Criterion struct {
	Family      string `yaml:"family"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	Satisfied   bool 
	SatisfiedBy []string
}

type Framework struct {
	Name     string             `yaml:"name"`
	Criteria map[string]Criterion `yaml:",inline"`
}

type Target map[string][]string

// CriteriaSatisfied determines the unique criteria currently satisfied by all Narratives, Policies, and Procedures
func CriteriaSatisfied(data *Data) map[string][]string {
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
		for _, criteriaKeys := range n.Satisfies {
			for _, key := range criteriaKeys {
				satisfied[key] = appendSatisfaction(satisfied, key, n.OutputFilename)
			}
		}
	}
	for _, n := range data.Policies {
		for _, criteriaKeys := range n.Satisfies {
			for _, key := range criteriaKeys {
				satisfied[key] = appendSatisfaction(satisfied, key, n.OutputFilename)
			}
		}
	}
	for _, n := range data.Controls {
		for _, criteriaKeys := range n.Targets {
			for _, key := range criteriaKeys {
				satisfied[key] = appendSatisfaction(satisfied, key, n.OutputFilename)
			}
		}
	}
	for _, n := range data.Procedures {
		for _, criteriaKeys := range n.Satisfies {
			for _, key := range criteriaKeys {
				satisfied[key] = appendSatisfaction(satisfied, key, n.OutputFilename)
			}
		}
	}
	return satisfied
}
