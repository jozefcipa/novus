package brew

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jozefcipa/novus/internal/logger"
)

var BrewPath string

var svcStartCommands []string
var svcStopCommands []string
var svcStatusCommands []string

type BrewServiceStatus struct {
	Running bool `json:"running"`
}

func init() {
	out, err := exec.Command("brew", "--prefix").Output()
	if err != nil {
		logger.Errorf("Failed to run \"brew --prefix\": %v", err)
		os.Exit(1)
	}

	BrewPath = strings.Replace(string(out), "\n", "", 1)
	svcStartCommands = []string{"brew", "services", "restart"}
	svcStopCommands = []string{"brew", "services", "stop"}
	svcStatusCommands = []string{"brew", "services", "info", "--json"}
}

func CheckIfRequiredBinariesInstalled() error {
	notInstalledMsg := "%s is not installed on this system!"

	if exists := binExists("nginx"); !exists {
		return fmt.Errorf(notInstalledMsg, "Nginx")
	}

	if exists := binExists("dnsmasq"); !exists {
		return fmt.Errorf(notInstalledMsg, "DNSMasq")
	}

	if exists := binExists("mkcert"); !exists {
		return fmt.Errorf(notInstalledMsg, "mkcert")
	}

	return nil
}

func InstallBinaries() error {
	// First check that Homebrew is installed
	brewExists := binExists("brew")
	if !brewExists {
		return &BrewMissingError{}
	}

	// Install required binaries - brew installs always latest by default
	// In case this causes problems in the future, we should consider pining to a specific version instead
	// e.g. https://cmichel.medium.com/how-to-install-an-old-package-version-with-brew-cc1c567dd088
	if exists := binExists("nginx"); !exists {
		if err := brewInstall("nginx"); err != nil {
			return err
		}
	}

	if exists := binExists("dnsmasq"); !exists {
		if err := brewInstall("dnsmasq"); err != nil {
			return err
		}
	}

	if exists := binExists("mkcert"); !exists {
		if err := brewInstall("mkcert"); err != nil {
			return err
		}
	}

	return nil
}

func RestartService(svc string) {
	cmds := append(svcStartCommands, svc)

	execBrewCommand(cmds)
}

func StopService(svc string) {
	cmds := append(svcStopCommands, svc)

	execBrewCommand(cmds)
}

func checkService(svc string, cmdOutput []byte) bool {
	// Parse json output
	var svcStatus []BrewServiceStatus
	json.Unmarshal(cmdOutput, &svcStatus)

	isRunning := len(svcStatus) > 0 && svcStatus[0].Running
	logger.Debugf("Service status of \"%s\" [running=%t]", svc, isRunning)

	return isRunning
}

func IsServiceRunning(svc string) bool {
	cmds := append(svcStatusCommands, svc)
	out := execBrewCommand(cmds)
	return checkService(svc, out)
}

func brewInstall(bin string) error {
	logger.Infof("⏳ Installing %s...", bin)
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
		return fmt.Errorf("An error occurred while installing \"%s\".\n\n%+v", bin, err)
	}
	fmt.Println() // print empty line

	// Check whether the binary is discoverable (in $PATH)
	binaryCheckCmd := exec.Command("which", bin)
	out, _ := binaryCheckCmd.Output()
	if string(out) == "" {
		logger.Errorf("%s has been installed but cannot be executed.", bin)
		logger.Hintf("The binary is probably not registered in the $PATH variable.")
		logger.Infof("   Run \"brew doctor\" or view https://github.com/jozefcipa/novus/issues/3 for more information.")
		os.Exit(1)
	} else {
		logger.Debugf("Binary '%s' available in %s", bin, out)
	}

	logger.Successf("%s installed", bin)

	return nil
}

func binExists(bin string) bool {
	_, err := exec.LookPath(bin)
	exists := err == nil

	logger.Debugf("Checking if binary exists [%s=%t]", bin, exists)

	return exists
}

func execBrewCommand(commands []string) []byte {
	commandString := strings.Join(commands, " ")
	logger.Debugf("Running \"%s\"", commandString)
	cmd := exec.Command(commands[0], commands[1:]...)

	out, err := cmd.Output()
	if err != nil {
		logger.Errorf("Failed to run \"%s\": %v", commandString, err)
		os.Exit(1)
	}

	return out
}
