package cmd

import (
	"fmt"
	"os"

	"github.com/arsham/figurine/figurine"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	// version is passed down from the main.go, this is only a placeholder to enable the --version flag
	Version: "-",
	Use:     "novus",
	Short:   "Local web development done effortlessly",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fs.ResolveDirs()
	},
	Run: func(cmd *cobra.Command, args []string) {
		figurine.Write(os.Stdout, "Novus", "ANSI Regular.flf")
		fmt.Println(`Novus is a tool that improves developer experience when working
on one or multiple web services by automatically providing
SSL secured URLs that proxy traffic to your services.

That means no more http://localhost:3000 calls.
Instead, open the "novus.yml" configuration and add a nice custom domain name
that will forward all the traffic to your upstream service.
To start run "novus serve --create-config" to initialize Novus and create an example configuration.`)

		cmd.Help()
	},
}

func Execute(version string) {
	// show app version
	rootCmd.SetVersionTemplate(version)

	// configure colors
	cc.Init(&cc.Config{
		RootCmd:  rootCmd,
		Headings: cc.HiCyan + cc.Bold + cc.Underline,
		Commands: cc.HiYellow + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Bold,
		Flags:    cc.Bold,
	})

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&logger.DebugEnabled, "debug", false, "include debug logs")
}
