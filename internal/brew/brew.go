package brew

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jozefcipa/novus/internal/logger"
)

var BrewPath string

var svcStartCommands []string
var svcStopCommands []string
func init() {
	out, err := exec.Command("brew", "--prefix").Output()
	if err != nil {
		logger.Errorf("Failed to run \"brew --prefix\": %v", err)
		os.Exit(1)
	}

	BrewPath = strings.Replace(string(out), "\n", "", 1)
	svcStartCommands = []string{"brew", "services", "restart"}
	svcStopCommands = []string{"brew", "services", "stop"}
}

func InstallBinaries() {
	// First check that Homebrew is installed
	brewExists := binExists("brew")
	if !brewExists {
		logger.Errorf("Novus requires Homebrew installed\n")
		logger.Infof("You can get it on https://brew.sh/\n\n.")
		os.Exit(1)
	}

	// Install required binaries
	if exists := binExists("nginx"); !exists {
		brewInstall("nginx@1.25")
	}

	if exists := binExists("dnsmasq"); !exists {
		brewInstall("dnsmasq@2.90")
	}

	if exists := binExists("mkcert"); !exists {
		brewInstall("mkcert@1.4")
	}
}

func RestartService(svc string) {
	cmds := append(svcStartCommands, svc)

	execBrewCommand(cmds)
}

func RestartServiceWithSudo(svc string) {
	// prepend with "sudo" and add "svc" at the end
	cmds := append([]string{"sudo"}, append(svcStartCommands, svc)...)

	execBrewCommand(cmds)
}

func StopService(svc string) {
	cmds := append(svcStopCommands, svc)

	execBrewCommand(cmds)
}

func StopServiceWithSudo(svc string) {
	// prepend with "sudo" and add "svc" at the end
	cmds := append([]string{"sudo"}, append(svcStopCommands, svc)...)

	execBrewCommand(cmds)
}

}

func brewInstall(bin string) {
	logger.Messagef("⏳ Installing %s...\n", bin)
	logger.Debugf("Running \"brew install %s\"", bin)

	cmd := exec.Command("brew", "install", bin)
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	err := cmd.Wait()

	if err != nil {
		logger.Errorf("An error occurred while installing \"%s\".\n\n%+v", bin, err)
		os.Exit(1)
	}

	logger.Successf("\n✅ %s installed\n", bin)
}

func binExists(bin string) bool {
	_, err := exec.LookPath(bin)
	exists := err == nil

	logger.Debugf("Checking if binary [%s] exists: %t", bin, exists)

	return exists
}

func execBrewCommand(commands []string) []byte {
	commandString := strings.Join(commands, " ")
	logger.Debugf("Running \"%s\"", commandString)
	cmd := exec.Command(commands[0], commands[1:]...)

	out, err := cmd.Output()
	if err != nil {
		logger.Errorf("Failed to run %s: %v", commandString, err)
		os.Exit(1)
	}

	return out
}
