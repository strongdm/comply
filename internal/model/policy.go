package model

type Policy struct {
  Name    string `yaml:"policyName"`
  ID      string `yaml:"policyID"`
  Clause  string `yaml:"policyClause"`
}
