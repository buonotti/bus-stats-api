package util

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

func ConfigValue(value string) string {
	return strings.ReplaceAll(value, "{env}", Env.String())
}
