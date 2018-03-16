package model

import "time"

type Procedure struct {
	Name       string       `yaml:"name"`
	Code       string       `yaml:"code"`
	Cron       string       `yaml:"cron"`
	Satisfies  Satisfaction `yaml:"satisfies"`
	FullPath   string
	ModifiedAt time.Time
	Body       string
}
