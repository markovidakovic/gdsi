// env.go handles the loading and parsing of env files.
// It provides functions to read key-value pairs from the files
// and set them in the applications environment
package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// loadEnv loads environment variables from a list of provided file paths.
// Each file is read in sequence, and its key-value pairs are set in the environment.
// If any file fails to load or contains errors, the process stops, and an error is returned.
func loadEnv(filenames []string) (err error) {
	for _, f := range filenames {
		envVars, err := readEnvFile(f)
		if err != nil {
			return err
		}
		err = setEnvVars(envVars)
		if err != nil {
			return err
		}
	}
	return
}

// readEnvFile reads an environment variables file and parses its key-value pairs
// into a map. The file should follow the format: KEY=VALUE, with optional quotes
// arount the value. Lines starting with '#' or empty lines are ignored.
func readEnvFile(filename string) (map[string]string, error) {
	envVars := make(map[string]string)

	// Open the file
	f, err := os.Open(filename)
	if err != nil {
		return envVars, fmt.Errorf("error loading environment file %s -> %w", filename, err)
	}
	defer f.Close()

	// Start scanning the file
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		key, val, err := parseEnvFileLine(scanner.Text())
		if err != nil {
			return envVars, fmt.Errorf("error parsing line in environment file %q -> %w", filename, err)
		}

		// Skip if key empty (line commented or empty)
		if key == "" {
			continue
		}

		envVars[key] = val
	}

	// Check for scanning error
	err = scanner.Err()
	if err != nil {
		return envVars, fmt.Errorf("error scanning the environment variables file %q -> %w", filename, err)
	}

	return envVars, nil
}

// parseEnvFileLine parses a single line from an environment variables file.
// The line should follow the format KEY=VALUE, with optional quotes around the value.
// Lines starting with '#' or empty lines are treated as comments and ignored.
func parseEnvFileLine(line string) (key string, value string, err error) {
	line = strings.TrimSpace(line)

	// Check if line is commented or is empty
	if strings.HasPrefix(line, "#") || len(line) == 0 {
		return "", "", nil
	}

	// Trim inline comments
	if commIdx := strings.Index(line, "#"); commIdx != -1 {
		line = strings.TrimSpace(line[:commIdx])
	}

	// Split key-value pairs
	lineParts := strings.SplitN(line, "=", 2)

	if len(lineParts) < 2 {
		return "", "", fmt.Errorf("invalid line %s", line)
	}

	key, value = strings.TrimSpace(lineParts[0]), strings.TrimSpace(lineParts[1])

	// Handle invalid key
	if !regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`).MatchString(key) {
		return "", "", fmt.Errorf("invalid environment variable key: %q", key)
	}

	// Handle empty value
	if value == "" {
		return "", "", fmt.Errorf("environment variable %q has an empty value", key)
	}

	// Remove quotes around values, if present
	if strings.HasPrefix(value, "\"") || strings.HasSuffix(value, "\"") {
		value = strings.Trim(value, "\"")
	} else if strings.HasPrefix(value, "'") || strings.HasSuffix(value, "'") {
		value = strings.Trim(value, "'")
	}

	return
}

// setEnvVars sets environment variables from a map of key-value pairs.
// If a variable already exists in the environment, it logs a warning and overwrites the value.
func setEnvVars(envVars map[string]string) error {
	for k, v := range envVars {
		// Check if var already exist
		if existing := os.Getenv(k); existing != "" {
			log.Printf("WARNING: overwriting existing environment variable %q: %s -> %s", k, existing, v)
		}

		// Set the env var
		err := os.Setenv(k, v)
		if err != nil {
			return fmt.Errorf("failed setting environment variable %q -> %w", k, err)
		}
		// log.Printf("setting environment variable: %s = %s", k, v)
	}

	return nil
}

// getEnvVar retrieves the value of an environment variable by its key.
// If the key does not exist or is empty, the specified defaultValue is returned.
func getEnvVar(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
