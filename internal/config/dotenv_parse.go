package config

import (
	"fmt"
	"os"
	"strings"
)

func applyDotEnvLine(line string) error {
	trimmed := strings.TrimSpace(line)
	if isIgnoredDotEnvLine(trimmed) {
		return nil
	}
	return parseDotEnvPair(trimmed)
}

func isIgnoredDotEnvLine(line string) bool {
	return line == "" || strings.HasPrefix(line, "#")
}

func parseDotEnvPair(line string) error {
	key, value, found := strings.Cut(line, "=")
	if !found {
		return nil
	}
	return setEnvIfNotSet(strings.TrimSpace(key), trimQuotes(strings.TrimSpace(value)))
}

func trimQuotes(value string) string {
	if len(value) < 2 {
		return value
	}
	if hasMatchingQuotes(value) {
		return value[1 : len(value)-1]
	}
	return value
}

func hasMatchingQuotes(value string) bool {
	return isQuoted(value, '"') || isQuoted(value, '\'')
}

func isQuoted(value string, quote rune) bool {
	return rune(value[0]) == quote && rune(value[len(value)-1]) == quote
}

func setEnvIfNotSet(key, value string) error {
	if key == "" {
		return nil
	}
	if _, exists := os.LookupEnv(key); exists {
		return nil
	}
	return setEnvValue(key, value)
}

func setEnvValue(key, value string) error {
	if err := os.Setenv(key, value); err != nil {
		return fmt.Errorf("set env %q: %w", key, err)
	}
	return nil
}
