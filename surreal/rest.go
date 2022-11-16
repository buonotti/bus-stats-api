package surreal

import (
	"strconv"

	"github.com/buonotti/bus-stats-api/config"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

var restClient = resty.New().
	SetBasicAuth("root", "root").
	SetHeader("Content-Type", "text/plain").
	SetHeader("NS", viper.GetString("database.{env}.ns")).
	SetHeader("DB", viper.GetString("database.{env}.db")).
	SetHeader("Accept", "application/json").
	SetDisableWarn(true)

func Url() string {
	protocol := viper.GetString(config.Get("database.{env}.protocol"))
	host := viper.GetString(config.Get("database.{env}.host"))
	port := viper.GetInt(config.Get("database.{env}.port"))

	return protocol + "://" + host + ":" + strconv.Itoa(port) + "/sql"
}
