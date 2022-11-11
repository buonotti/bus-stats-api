package util

import (
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

var RestClient = resty.New().
	SetBasicAuth("root", "root").
	SetHeader("Content-Type", "text/plain").
	SetHeader("NS", viper.GetString("database.{env}.ns")).
	SetHeader("DB", viper.GetString("database.{env}.db")).
	SetHeader("Accept", "application/json").
	SetDisableWarn(true)
