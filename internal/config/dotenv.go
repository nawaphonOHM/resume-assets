package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

var defaultDotEnvFiles = []string{".env", "configs/.env"}

func loadDotEnvFiles(filenames ...string) error {
	if len(filenames) == 0 {
		return loadDefaultDotEnvFiles()
	}
	return loadExplicitDotEnvFiles(filenames)
}

func loadDefaultDotEnvFiles() error {
	for _, filename := range defaultDotEnvFiles {
		err := loadDotEnvFile(filename)
		if shouldStopLoadingDefault(err) {
			return err
		}
	}
	return nil
}

func shouldStopLoadingDefault(err error) bool {
	return err == nil || !errors.Is(err, os.ErrNotExist)
}

func loadExplicitDotEnvFiles(filenames []string) error {
	for _, filename := range filenames {
		if err := loadDotEnvFile(filename); err != nil {
			return err
		}
	}
	return nil
}

func loadDotEnvFile(filename string) (result error) {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open env file: %w", err)
	}
	defer closeDotEnvFile(file, &result)

	return scanDotEnvFile(file)
}

func closeDotEnvFile(file *os.File, result *error) {
	if err := file.Close(); err != nil {
		*result = errors.Join(*result, fmt.Errorf("close env file: %w", err))
	}
}

func scanDotEnvFile(file *os.File) error {
	scanner := bufio.NewScanner(file)
	if err := scanDotEnvLines(scanner); err != nil {
		return err
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan env file: %w", err)
	}
	return nil
}

func scanDotEnvLines(scanner *bufio.Scanner) error {
	for scanner.Scan() {
		if err := applyDotEnvLine(scanner.Text()); err != nil {
			return err
		}
	}
	return nil
}
