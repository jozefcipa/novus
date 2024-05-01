package tui

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/novus"
)

func AskUser(prompt string) string {
	logger.Messagef(prompt)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
		logger.Errorf("Failed to read from CLI: %v", err)
		os.Exit(1)
	}

	return scanner.Text()
}

func Confirm(question string) bool {
	// TODO: implement
	return false
}

func PrintRoutingTable(apps map[string]*novus.AppState) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"Application", "Upstream ", "Domain", "Status", "Directory"})

	for appName, appState := range apps {
		for _, route := range appState.Routes {
			t.AppendRow(table.Row{appName, route.Upstream, fmt.Sprintf("https://%s", route.Domain), "🚀", appState.Directory}, table.RowConfig{AutoMerge: true})
		}
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true, Colors: text.Colors{text.FgCyan}},
		{Number: 2, AutoMerge: false},
		{Number: 3, AutoMerge: false, Colors: text.Colors{text.FgHiGreen}},
		{Number: 4, AutoMerge: false, Align: text.AlignCenter}, // TODO: set to true when figure out how to vertical align to middle
		{Number: 5, AutoMerge: true},
	})

	t.SortBy([]table.SortBy{{Name: "Application", Mode: table.Asc}})
	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = true

	t.Render()
}
