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
	fmt.Print(prompt)

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
	allApps := novusState.GetAllApps()
	if len(allApps) == 0 {
		logger.Warnf("You don't have any apps configured.")
		logger.Hintf("Run \"novus init\" or \"novus serve\" to start routing.")
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
		appState := allApps[appName]
		color := tablewriter.FgGreenColor
		if appState.Status == novus.APP_PAUSED {
			color = tablewriter.FgYellowColor
		}

		for _, route := range appState.Routes {
			table.Rich(
				[]string{
					appName,
					route.Upstream,
					fmt.Sprintf("https://%s", route.Domain),
					strings.ToUpper(string(appState.Status)),
					appState.Directory,
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
}

func ParseAppFromArgs(args []string) (string, *novus.AppState) {
	if len(args) < 1 {
		logger.Errorf("App name not provided!")
		logger.Hintf("Please specify app that you want to remove by running \"novus remove [app-name]\"")
		os.Exit(1)
	}
	appName := args[0]

	// Load app state for the given app name if it exists, or throw an error
	appState, exists := novus.GetAppState(appName)
	if !exists {
		logger.Errorf("App name \"%s\" is not registered in Novus", appName)
		os.Exit(1)
	}

	return appName, appState
}
