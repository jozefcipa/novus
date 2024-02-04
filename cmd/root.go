package cmd

import (
	"os"

	"github.com/jozefcipa/novus/internal/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "novus",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&logger.DebugEnabled, "debug", false, "include debug logs")
}
