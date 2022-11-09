package util

import (
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

var RestClient = resty.New().
	SetBasicAuth(viper.GetString(GetConfig("database.{env}.user")), viper.GetString(GetConfig("database.{env}.pass"))).
	SetHeader("Content-Type", "application/json").
	SetHeader("NS", viper.GetString("database.{env}.ns")).
	SetHeader("DB", viper.GetString("database.{env}.db")).
	SetDisableWarn(true)
