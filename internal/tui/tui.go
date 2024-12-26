package tui

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/shared"
	"github.com/olekukonko/tablewriter"
)

func AskUser(prompt string) string {
	fmt.Printf("%s%s%s", logger.GRAY, prompt, logger.RESET)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	err := scanner.Err()
	if err != nil {
		logger.Errorf("Failed to read from CLI: %v", err)
		os.Exit(1)
	}

	return scanner.Text()
}

func Confirm(question string) bool {
	answer := AskUser(fmt.Sprintf("%s [Y/n]: ", question))
	return answer == "Y"
}

func PrintRoutingTable(novusState novus.NovusState) {
	allApps := novusState.Apps
	hasSomeRoutes := false
	for appName, appState := range allApps {
		if appName != novus.NovusInternalAppName && len(appState.Routes) > 0 {
			hasSomeRoutes = true
			break
		}
	}
	if !hasSomeRoutes {
		logger.Warnf("You don't have any apps configured.")
		logger.Hintf(" Run \"novus init\" to configure routing.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Application", "Upstream ", "Domain", "Status", "Directory"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: true})
	table.SetAutoMergeCellsByColumnIndex([]int{0, 3, 4})
	table.SetCenterSeparator("|")
	table.SetRowLine(true)

	// Sort apps
	sortedAppNames := shared.MapKeys(allApps)
	slices.SortFunc(sortedAppNames, func(a, b string) int { return cmp.Compare(a, b) })

	for _, appName := range sortedAppNames {
		// Only print out internal routes in debug mode
		if appName == novus.NovusInternalAppName && !logger.DebugEnabled {
			continue
		}

		appState := allApps[appName]
		color := tablewriter.FgGreenColor
		if appState.Status == novus.APP_PAUSED {
			color = tablewriter.FgYellowColor
		}

		for _, route := range appState.Routes {
			displayAppName := appName
			displayDir := appState.Directory
			if appName == novus.GlobalAppName {
				displayAppName = "Global Routes"
				displayDir = ""
			}

			table.Rich(
				[]string{
					displayAppName,
					route.Upstream,
					fmt.Sprintf("https://%s", route.Domain),
					strings.ToUpper(string(appState.Status)),
					displayDir,
				},
				[]tablewriter.Colors{
					{tablewriter.FgCyanColor},
					{tablewriter.UnderlineSingle},
					{tablewriter.Bold, tablewriter.UnderlineSingle, color},
					{color},
					{},
				})
		}
	}

	fmt.Println() // print empty line
	table.Render()
	fmt.Println()
	logger.Hintf("You can also view these routes in your browser at %shttps://index.novus%s", logger.UNDERLINE, logger.RESET)
}

func ParseAppFromArgs(args []string, cmd string) (string, *novus.AppState) {
	if len(args) < 1 {
		logger.Errorf("App name not provided!")
		logger.Hintf("Please specify app name by running \"novus %s [app-name]\"", cmd)
		os.Exit(1)
	}
	appName := args[0]

	if appName == novus.NovusInternalAppName || appName == novus.GlobalAppName {
		logger.Errorf("App \"%s\" is used by Novus and cannot be %sd", appName, cmd)
		os.Exit(1)
	}

	// Load app state for the given app name if it exists, or throw an error
	appState, exists := novus.GetAppState(appName)
	if !exists {
		return appName, nil
	}

	return appName, appState
}
