package cli

import (
	"bufio"
	"log"
	"os"

	"github.com/jozefcipa/novus/internal/logger"
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
