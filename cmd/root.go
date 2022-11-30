package cmd

import (
	"github.com/buonotti/bus-stats-api/errors"
	"github.com/spf13/cobra"
)

// Root command
var rootCmd = &cobra.Command{
	Version: "0.0.1-beta.2",
	Use:     "bus-stats-api",
	Short:   "Bus Stats API",
	Long:    `Api providing the services for the bus-stats application`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("debug").Value.String() == "true" {
			errors.DoDebug = true
		}
	},
}

// Set the version template
func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")
	rootCmd.SetVersionTemplate("{{.Version}}\n")
}

// Execute the root command
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
