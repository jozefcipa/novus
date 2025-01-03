package domain_cleanup_manager

import (
	"github.com/jozefcipa/novus/internal/diff_manager"
	"github.com/jozefcipa/novus/internal/dns_manager"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/sharedtypes"
	"github.com/jozefcipa/novus/internal/ssl_manager"
)

func RemoveDomains(routes []sharedtypes.Route, appName string, novusState *novus.NovusState) {
	appState, _ := novus.GetAppState(appName)

	if len(routes) == 1 {
		ssl_manager.DeleteCert(routes[0].Domain, appState)
		logger.Checkf("Removed domain [%s]", routes[0].Domain)
	} else {
		logger.Checkf("Removed %d domains:", len(routes))
		for _, deletedRoute := range routes {
			logger.Infof("   - %s", deletedRoute.Domain)
			ssl_manager.DeleteCert(deletedRoute.Domain, appState)
		}
	}

	// Remove DNS records for unused TLDs
	otherAppsRoutes := []sharedtypes.Route{}
	for novusAppName, novusAppState := range novusState.GetActiveApps() {
		// We want to find usage only in other apps,
		// current app's state has not yet been updated so it contains routes that we're deleting,
		// thus this would yield false results claiming the TLD is still used
		if novusAppName != appName {
			otherAppsRoutes = append(otherAppsRoutes, novusAppState.Routes...)
		}
	}

	unusedTLDs := diff_manager.DetectUnusedTLDs(routes, otherAppsRoutes)
	if len(unusedTLDs) > 0 {
		for _, tld := range unusedTLDs {
			logger.Debugf("Removing unused TLD domain [*.%s]", tld)
			dns_manager.UnregisterTLD(tld, novusState)
		}
	}
}
