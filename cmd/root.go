package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/arsham/figurine/figurine"
	"github.com/fatih/color"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/paths"
	"github.com/jozefcipa/novus/internal/sharedtypes"
	"github.com/jozefcipa/novus/internal/tld"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	// Version is passed down from the main.go,
	// this is only a placeholder to enable the --version flag
	Version: "-",
	Use:     "novus",
	Short:   "Local web development done effortlessly",

	// Here is the init code that runs before any command
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		paths.Resolve()
		tld.LoadExistingTLDsFile()
		novusState := novus.GetState()

		ctx := cmd.Context().Value(sharedtypes.CommandContext{}).(sharedtypes.CommandContext)

		if novusState.Version != ctx.Version {
			// logger.Infof("🚀 New Novus version detected [%s]. Updating state...", ctx.Version)
			// Here we can update the configuration/state in the future if needed after a new version is released

			novusState.Version = ctx.Version
			novus.SaveState()
		}
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

func Execute(version string, buildDate string) {
	// inject version into command context
	ctx := context.WithValue(context.TODO(), sharedtypes.CommandContext{}, sharedtypes.CommandContext{
		Version: version,
	})
	rootCmd.SetContext(ctx)

	// Set app version
	versionTemplate := fmt.Sprintf("Novus v%s (built on %s) %s/%s\n", version, buildDate, runtime.GOOS, runtime.GOARCH)
	rootCmd.SetVersionTemplate(versionTemplate)

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
