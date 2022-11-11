package cmd

import (
	"github.com/spf13/cobra"
)

// Root command
var rootCmd = &cobra.Command{
	Version: "0.0.1-beta.2",
	Use:     "bus-stats-api",
}

func init() {
	rootCmd.SetVersionTemplate("{{.Version}}\n")
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
