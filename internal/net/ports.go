package net

import (
	"fmt"
	go_net "net"
	"os"
	"os/exec"
	"strings"

	"github.com/jozefcipa/novus/internal/logger"
)

func lsof(ports []string) []string {
	commandString := fmt.Sprintf("sudo lsof -nP -i4:%s", strings.Join(ports, ","))
	logger.Debugf("Running \"%s\"", commandString)

	cmd := exec.Command("sudo", "lsof", "-nP", fmt.Sprintf("-i4:%s", strings.Join(ports, ",")))
	out, err := cmd.CombinedOutput()
	result := string(out)

	// `lsof` command returns exit code 1 even if there is no script error but no information was found
	// therefore, let's only return error if the exit code is 1 and there is some actual output
	// https://stackoverflow.com/a/29843137/4480179
	if err != nil && result != "" {
		logger.Errorf("Failed to run \"%s\": %v\n%s", commandString, err, result)
		os.Exit(1)
	}

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
		_, port, _ := go_net.SplitHostPort(recordParts[8])

		portUsage[port] = binary
	}

	logger.Debugf("Port usage: %v", portUsage)

	return portUsage
}

type PortUsage = map[string]string

func CheckPortsUsage(ports ...string) PortUsage {
	lsof := lsof(ports)
	logger.Debugf("lsof result:\n%s", strings.Join(lsof, "\n"))

	return parseLsof(lsof)
}
