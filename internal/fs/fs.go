package fs

import (
	"os"

	"github.com/jozefcipa/novus/internal/logger"
)

func ReadFileOrExit(path string) string {
	file, err := os.ReadFile(path)
	if err != nil {
		logger.Errorf("Failed to read a file %s\n   Reason: %v", path, err)
		os.Exit(1)
	}

	return string(file)
}

func ReadFile(path string) (string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(file), nil
}

func WriteFileOrExit(path string, data string) {
	err := os.WriteFile(path, []byte(data), 0644)
	if err != nil {
		logger.Errorf("Failed to write to a file %s\n   Reason: %v", path, err)
		os.Exit(1)
	}
}

func DeleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		logger.Errorf("Failed to delete file %s\n   Reason: %v", path, err)
		return err
	}
	return nil
}

func MakeDirOrExit(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		logger.Errorf("Failed to create directory %s\n   Reason: %v", path, err)
		os.Exit(1)
	}
}

func DeleteDir(path string) error {
	if err := os.RemoveAll(path); err != nil {
		logger.Errorf("Failed to delete directory %s\n   Reason: %v", path, err)
		return err
	}
	return nil
}

func FileExists(path string) bool {
	fInfo, _ := os.Stat(path)

	return fInfo != nil
}
