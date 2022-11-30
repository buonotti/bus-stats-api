package config

import (
	"github.com/buonotti/bus-stats-api/errors"
	"github.com/spf13/viper"
)

func Setup() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		err = errors.ConfigFileParseError.WrapWithNoMessage(err)
	}
	errors.CheckError(err)
}
