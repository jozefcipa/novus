package diff_manager

import (
	"slices"

	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/dns_manager"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/shared"
)

func routeExists(domain string, routes []shared.Route) bool {
	for _, route := range routes {
		if route.Domain == domain {
			return true
		}
	}

	return false
}

func DetectConfigDiff(conf config.NovusConfig, state novus.AppState) (added []shared.Route, deleted []shared.Route) {
	// Detect routes that are stored in state but have been removed from the configuration file
	for _, route := range state.Routes {
		if !routeExists(route.Domain, conf.Routes) {
			deleted = append(deleted, route)
		}
	}

	// Detect routes that are found in configuration file but are not stored in the state
	for _, route := range conf.Routes {
		if !routeExists(route.Domain, state.Routes) {
			added = append(added, route)
		}
	}

	return added, deleted
}

func DetectUnusedTLDs(deletedRoutes []shared.Route, stateRoutes []shared.Route) []string {
	deletedRoutesTLDs := dns_manager.GetTLDs(deletedRoutes)
	stateTLDs := dns_manager.GetTLDs(stateRoutes)
	unusedTLDs := []string{}

	// Iterate through all routes that have been deleted
	// and check if their TLD domain is used in the remaining routes in state
	// if not, that means the TLD is not used anymore and can be removed
	for _, deletedRouteTLD := range deletedRoutesTLDs {
		if !slices.Contains(stateTLDs, deletedRouteTLD) {
			unusedTLDs = append(unusedTLDs, deletedRouteTLD)
		}
	}

	return unusedTLDs
}

type appDomain struct {
	App    string
	Domain string
}

func DetectDuplicateDomains(existingApps map[string]novus.AppState, addedRoutes []shared.Route) error {
	allDomains := []appDomain{}

	// Collect all existing domains across apps
	for appName, appConfig := range existingApps {
		for _, route := range appConfig.Routes {
			allDomains = append(allDomains, appDomain{App: appName, Domain: route.Domain})
		}
	}

	// Iterate through the newly added routes to see if some of them already exists in the slice
	for _, route := range addedRoutes {
		if idx := slices.IndexFunc(allDomains, func(appDomain appDomain) bool { return appDomain.Domain == route.Domain }); idx != -1 {
			return &DuplicateDomainError{
				DuplicateDomain:       route.Domain,
				OriginalAppWithDomain: allDomains[idx].App,
			}
		}
	}

	return nil
}
