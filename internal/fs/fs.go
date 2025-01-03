package fs

import (
	"io"
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

func WriteFileWithSudoOrExit(path string, data string) {
	if _, err := exec.Command("sudo", "touch", path).Output(); err != nil {
		logger.Errorf("Failed to create file %s\n   Reason: %v", path, err)
		os.Exit(1)
	}

	// We need to change the file owner to the current user in order to be able to write to the file
	user, _ := user.Current()
	ChownOrExit(path, user.Username)
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

func Copy(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func ChownOrExit(path string, user string) {
	if _, err := exec.Command("sudo", "chown", user, path).Output(); err != nil {
		logger.Errorf("Failed to call `chown` on file %s\n   Reason: %v", path, err)
		os.Exit(1)
	}
}
