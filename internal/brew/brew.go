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

func init() {
	out, err := exec.Command("brew", "--prefix").Output()
	if err != nil {
		logger.Errorf("Failed to run \"brew --prefix\": %v", err)
		os.Exit(1)
	}

	BrewPath = strings.Replace(string(out), "\n", "", 1)
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
		brewInstall("nginx")
	}

	if exists := binExists("dnsmasq"); !exists {
		brewInstall("dnsmasq")
	}

	if exists := binExists("mkcert"); !exists {
		brewInstall("mkcert")
	}
}

func StartBrewService(svc string) {
	logger.Debugf("Running \"brew services start %s\"", svc)
	cmd := exec.Command("brew", "services", "start", svc)

	err := cmd.Run()
	if err != nil {
		logger.Errorf("Failed to start %s: %v", svc, err)
		os.Exit(1)
	}
}

func RestartBrewServiceWithSudo(svc string) {
	logger.Debugf("Running \"sudo brew services restart %s\"", svc)
	cmd := exec.Command("sudo", "brew", "services", "restart", svc)

	err := cmd.Run()
	if err != nil {
		logger.Errorf("Failed to restart %s: %v", svc, err)
		os.Exit(1)
	}
}

func StopBrewService(svc string) {
	logger.Debugf("Running \"brew services stop %s\"", svc)
	cmd := exec.Command("brew", "services", "stop", svc)

	err := cmd.Run()
	if err != nil {
		logger.Errorf("Failed to stop %s: %v", svc, err)
		os.Exit(1)
	}
}

func RestartBrewService(svc string) {
	logger.Debugf("Running \"brew services restart %s\"", svc)
	cmd := exec.Command("brew", "services", "restart", svc)

	err := cmd.Run()
	if err != nil {
		logger.Errorf("Failed to restart %s: %v", svc, err)
		os.Exit(1)
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
