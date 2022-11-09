package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Root command
var rootCmd = &cobra.Command{
	Version: viper.GetString("cli.version"),
	Use:     "bus-stats-api",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
