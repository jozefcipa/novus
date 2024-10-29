package cmd

import (
	"fmt"
	"os"

	"github.com/arsham/figurine/figurine"
	"github.com/fatih/color"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/net"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/ssl_manager"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	// Version is passed down from the main.go, this is only a placeholder to enable the --version flag
	Version: "-",
	Use:     "novus",
	Short:   "Local web development done effortlessly",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fs.ResolveDirs()
		novus.ResolveDirs()
		ssl_manager.ResolveDirs()
		net.LoadExistingTLDsFile()
	},
	Run: func(cmd *cobra.Command, args []string) {
		figurine.Write(os.Stdout, "Novus", "ANSI Regular.flf")

		white := color.New(color.FgHiWhite)
		whiteBold := white.Add(color.Bold)
		white.Print("Novus")
		fmt.Println(` is a tool that improves developer experience
by automatically provisioning user-friendly HTTPS URLs that proxy traffic to your services.
That means no more messing around with ` + color.New(color.Underline).Sprint("http://localhost:3000") + ` URLs 🙌

Start by running "` + whiteBold.Sprint("novus init") + `" to create a configuration file 🚀

Open the ` + whiteBold.Sprint(config.ConfigFileName) + ` configuration and define custom domain names
that will forward all the traffic to your upstream services.`)

		// Override the default HelpFunc to display only the "Available Commands" section
		title := color.New(color.FgCyan).Add(color.Underline)
		title.Print("\nAvailable Commands:\n\n")

		yellow := color.New(color.FgYellow)
		for _, c := range cmd.Commands() {
			if !c.Hidden {
				yellow.Printf("  %-15s ", c.Name())
				fmt.Printf("%s\n", c.Short)
			}
		}

		fmt.Print("\nUse \"")
		whiteBold.Print("novus ")
		fmt.Print("[command] --help\" for more information about a command.\n")
	},
}

func Execute(version string) {
	// Show app version
	rootCmd.SetVersionTemplate(version)

	// Configure colors
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
