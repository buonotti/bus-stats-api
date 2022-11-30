package env

import "strings"

type environment string

const (
	Development environment = "development"
	Production  environment = "production"
)

var Env environment

func (e environment) String() string {
	return string(e)
}

func Get(key string) string {
	return strings.ReplaceAll(key, "{env}", Env.String())
}
