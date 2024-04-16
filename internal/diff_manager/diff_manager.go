package diff_manager

import (
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/shared"
)

type DiffType int

const (
	Added DiffType = iota
	Deleted
)

type Diff struct {
	Route shared.Route
	Diff  DiffType
}

func routeExists(domain string, routes []shared.Route) bool {
	for _, route := range routes {
		if route.Domain == domain {
			return true
		}
	}

	return false
}

func DetectConfigDiff(conf config.NovusConfig, state novus.AppState) (added []Diff, deleted []Diff) {
	// detect routes that are stored in state but have been removed from the configuration file
	for _, route := range state.Routes {
		if !routeExists(route.Domain, conf.Routes) {
			deleted = append(deleted, Diff{
				Route: route,
				Diff:  Deleted,
			})
		}
	}

	// detect routes that are found in configuration file but are not stored in the state
	for _, route := range conf.Routes {
		if !routeExists(route.Domain, state.Routes) {
			added = append(added, Diff{
				Route: route,
				Diff:  Added,
			})
		}
	}

	return added, deleted
}

func DetectUnusedTLDs(conf config.NovusConfig, state novus.AppState) (unusedTLDs []string) {
	configTLDs := dnsmasq.GetTLDs(conf.Routes)
	stateTLDs := dnsmasq.GetTLDs(state.Routes)

	// get TLDs that exist in the state but not in the config
	unusedTLDs = shared.Difference(stateTLDs, configTLDs)

	return unusedTLDs
}
