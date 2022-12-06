package surreal

import (
	"strconv"

	"github.com/buonotti/bus-stats-api/config/env"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

func restClient() *resty.Client {
	return resty.New().
	SetBasicAuth(viper.GetString(env.Get("database.{env}.user")), viper.GetString(env.Get("database.{env}.pass"))).
	SetHeader("Content-Type", "text/plain").
	SetHeader("NS", viper.GetString(env.Get("database.{env}.ns"))).
	SetHeader("DB", viper.GetString(env.Get("database.{env}.db"))).
	SetHeader("Accept", "application/json").
	SetDisableWarn(true)
}

func Url() string {
	protocol := viper.GetString(env.Get("database.{env}.protocol"))
	host := viper.GetString(env.Get("database.{env}.host"))
	port := viper.GetInt(env.Get("database.{env}.port"))

	return protocol + "://" + host + ":" + strconv.Itoa(port) + "/sql"
}
