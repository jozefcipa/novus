package cmd

import (
	"os"

	"github.com/jozefcipa/novus/internal/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "novus",
	Short: "Local web development done effortlessly",
	Long: ` _   _
| \ | | _____   ___   _ ___
|  \| |/ _ \ \ / / | | / __|
| |\  | (_) \ V /| |_| \__ \
|_| \_|\___/ \_/  \__,_|___/

Novus is a tool that improves developer experience when working
on one or multiple web services by automatically providing
SSL secured URLs that proxy traffic to your services.

No more http://localhost:3000 calls.
Instead, open the "novus.yml" configuration and add a nice custom domain name
that will forward all the traffic to your upstream service.

To start run "novus serve --create-config" to initialize Novus and create an example configuration.
`,
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
