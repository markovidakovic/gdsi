package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

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

func readEnvFile(filename string) (map[string]string, error) {
	envVars := make(map[string]string)

	// open the file
	f, err := os.Open(filename)
	if err != nil {
		return envVars, fmt.Errorf("error loading environment file %s -> %w", filename, err)
	}
	defer f.Close()

	// start scanning the file
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		key, val, err := parseEnvFileLine(scanner.Text())
		if err != nil {
			return envVars, fmt.Errorf("error parsing line in environment file %q -> %w", filename, err)
		}

		// skip if key empty (line commented or empty)
		if key == "" {
			continue
		}

		envVars[key] = val
	}

	// check for scanning error
	err = scanner.Err()
	if err != nil {
		return envVars, fmt.Errorf("error scanning the environment variables file %q -> %w", filename, err)
	}

	return envVars, nil
}

func parseEnvFileLine(line string) (key string, value string, err error) {
	line = strings.TrimSpace(line)

	// check if line is commented or is empty
	if strings.HasPrefix(line, "#") || len(line) == 0 {
		return "", "", nil
	}

	// split key-value pairs
	lineParts := strings.SplitN(line, "=", 2)

	if len(lineParts) < 2 {
		return "", "", fmt.Errorf("invalid line %s", line)
	}

	key, value = strings.TrimSpace(lineParts[0]), strings.TrimSpace(lineParts[1])

	// handle invalid key
	if !regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`).MatchString(key) {
		return "", "", fmt.Errorf("invalid environment variable key: %q", key)
	}

	// handle invalid value
	if strings.Contains(value, "\n") {
		return "", "", fmt.Errorf("environment variable %q has an invalid multiline value", key)
	}

	// handle empty value
	if value == "" {
		return "", "", fmt.Errorf("environment variable %q has an empty value", key)
	}

	// remove quotes around values, if present
	if strings.HasPrefix(value, "\"") || strings.HasSuffix(value, "\"") {
		value = strings.Trim(value, "\"")
	} else if strings.HasPrefix(value, "'") || strings.HasSuffix(value, "'") {
		value = strings.Trim(value, "'")
	} else {
		// if not qouted, trim inline comments
		if commIdx := strings.Index(value, "#"); commIdx != -1 {
			value = strings.TrimSpace(value[:commIdx])
		}
	}

	return
}

func setEnvVars(envVars map[string]string) error {
	for k, v := range envVars {
		// check if var already exist
		if existing := os.Getenv(k); existing != "" {
			log.Printf("WARNING: overwriting existing environment variable %q: %s -> %s", k, existing, v)
		}

		// set the env var
		err := os.Setenv(k, v)
		if err != nil {
			return fmt.Errorf("failed setting environment variable %q -> %w", k, err)
		}
		// log.Printf("setting environment variable: %s = %s", k, v)
	}

	return nil
}

func getEnvVar(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
