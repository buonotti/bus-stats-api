package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Setup() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	cobra.CheckErr(err)
}
