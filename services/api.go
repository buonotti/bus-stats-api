package services

import (
	"strconv"

	"github.com/spf13/viper"
)

func ApiUrl(path string) string {
	protocol := viper.GetString(ConfigValue("database.{env}.protocol"))
	host := viper.GetString(ConfigValue("database.{env}.host"))
	port := viper.GetInt(ConfigValue("database.{env}.port"))

	return protocol + "://" + host + ":" + strconv.Itoa(port) + "/" + path
}
