package ports

import (
	"net"
	"strings"

	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/sudo"
)

func lsof(ports []string) []string {
	result := sudo.CheckPortsOrExit(ports)

	return strings.Split(
		strings.TrimRight(result, "\n"), // remove \n from the end of the string so we don't create an empty record in the array
		"\n",
	)[1:] // first line is header, so we skip it
}

// Example record format: "dnsmasq   31695    nobody    5u  IPv4 0x51fe684ad72f7a85      0t0  TCP 192.168.64.1:53 (LISTEN)"
func parseLsof(lsofRecords []string) map[string]string {
	portUsage := make(map[string]string)

	for _, record := range lsofRecords {
		// [0] - binary name
		// [8] - listen address
		// [9] - connection status (OPTIONAL)
		recordParts := strings.Fields(record)

		binary := recordParts[0]
		_, port, _ := net.SplitHostPort(recordParts[8])

		portUsage[port] = binary
	}

	logger.Debugf("Port usage: %v", portUsage)

	return portUsage
}

type PortUsage = map[string]string

func CheckIfAvailable(ports ...string) PortUsage {
	logger.Infof("Checking ports availability")
	lsof := lsof(ports)
	logger.Debugf("lsof result:\n%s", strings.Join(lsof, "\n"))

	return parseLsof(lsof)
}
