package cmd

import (
	"github.com/spf13/cobra"
)

// Root command
var rootCmd = &cobra.Command{
	Version: "0.0.1-beta.2",
	Use:     "bus-stats-api",
	Short:   "Bus Stats API",
	Long:    `Api providing the services for the bus-stats application`,
}

// Set the version template
func init() {
	rootCmd.SetVersionTemplate("{{.Version}}\n")
}

// Execute the root command
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
