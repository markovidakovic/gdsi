package config

import (
	"os"
	"testing"
)

func TestGetEnvVar(t *testing.T) {
	// Set environment variables
	os.Setenv("EXISTING_VAR", "test_value")
	defer os.Unsetenv("EXISTING_VAR")
	os.Setenv("EMPTY_VAR", "")
	defer os.Unsetenv("EMPTY_VAR")

	// Test data
	testCases := []struct {
		name         string
		key          string
		defaultValue string
		expected     string
	}{
		{
			name:         "ExistingVariable",
			key:          "EXISTING_VAR",
			defaultValue: "default_value",
			expected:     "test_value",
		},
		{
			name:         "NonExistingVariable",
			key:          "NON_EXISTING_VAR",
			defaultValue: "default_value",
			expected:     "default_value",
		},
		{
			name:         "EmptyVariable",
			key:          "EMPTY_VAR",
			defaultValue: "default_value",
			expected:     "default_value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getEnvVar(tc.key, tc.defaultValue)
			if result != tc.expected {
				t.Errorf("getEnvVar(%q, %q) = %q; want %q", tc.key, tc.defaultValue, result, tc.expected)
			}
		})
	}
}

func TestSetEnvVars(t *testing.T) {
	// Test data
	testCases := []struct {
		name          string
		input         map[string]string
		existingVars  map[string]string
		expectedVars  map[string]string
		expectedError bool
	}{
		{
			name: "SetNewEnvironmentVariables",
			input: map[string]string{
				"NEW_VAR1": "value1",
				"NEW_VAR2": "value2",
			},
			existingVars: nil,
			expectedVars: map[string]string{
				"NEW_VAR1": "value1",
				"NEW_VAR2": "value2",
			},
			expectedError: false,
		},
		{
			name: "OverwriteExistingEnvironmentVariables",
			input: map[string]string{
				"EXISTING_VAR": "new_value",
			},
			existingVars: map[string]string{
				"EXISTING_VAR": "old_value",
			},
			expectedVars: map[string]string{
				"EXISTING_VAR": "new_value",
			},
			expectedError: false,
		},
		{
			name:  "HandleEmptyMap",
			input: map[string]string{},
			existingVars: map[string]string{
				"EXISTING_VAR": "value",
			},
			expectedVars: map[string]string{
				"EXISTING_VAR": "value",
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set existing environment variables
			for k, v := range tc.existingVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			err := setEnvVars(tc.input)

			// Check for unexpected errors
			if (err != nil) != tc.expectedError {
				t.Errorf("setEnvVars() error = %v, expectedError %v", err, tc.expectedError)
			}

			// Verify the environment variables
			for k, expected := range tc.expectedVars {
				actual := os.Getenv(k)
				if actual != expected {
					t.Errorf("environment variable %q = %q, want %q", k, actual, expected)
				}
			}

			// Verify no unexpected variables are set
			for k := range tc.input {
				if _, exists := tc.expectedVars[k]; !exists {
					actual := os.Getenv(k)
					if actual != "" {
						t.Errorf("unexpected environment variable %q = %q", k, actual)
					}
				}
			}
		})
	}
}

func TestParseEnvFileLine(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedKey   string
		expectedValue string
		expectedError bool
	}{
		{
			name:          "ValidKeyValuePair",
			input:         "KEY=VALUE",
			expectedKey:   "KEY",
			expectedValue: "VALUE",
			expectedError: false,
		},
		{
			name:          "KeyValuePairWithSpaces",
			input:         " KEY = VALUE ",
			expectedKey:   "KEY",
			expectedValue: "VALUE",
			expectedError: false,
		},
		{
			name:          "KeyWithInvalidCharacters",
			input:         "INVALID-KEY=VALUE",
			expectedKey:   "",
			expectedValue: "",
			expectedError: true,
		},
		{
			name:          "KeyWithEmptyValue",
			input:         "KEY=",
			expectedKey:   "",
			expectedValue: "",
			expectedError: true,
		},
		{
			name:          "EmptyLine",
			input:         "  ",
			expectedKey:   "",
			expectedValue: "",
			expectedError: false,
		},
		{
			name:          "CommentedLine",
			input:         "# KEY=VALUE",
			expectedKey:   "",
			expectedValue: "",
			expectedError: false,
		},
		{
			name:          "InlineCommentLine",
			input:         "KEY=VALUE # this is a inline comment",
			expectedKey:   "KEY",
			expectedValue: "VALUE",
			expectedError: false,
		},
		{
			name:          "KeyValuePairWithDoubleQoutedValue",
			input:         "KEY=\"VALUE\"",
			expectedKey:   "KEY",
			expectedValue: "VALUE",
			expectedError: false,
		},
		{
			name:          "KeyValuePairWIthSingleQoutedValue",
			input:         "KEY='VALUE'",
			expectedKey:   "KEY",
			expectedValue: "VALUE",
			expectedError: false,
		},
		{
			name:          "QoutedEmptyValue",
			input:         "KEY=\"\"",
			expectedKey:   "KEY",
			expectedValue: "",
			expectedError: false,
		},
		{
			name:          "InvalidLineWithoutEqualsSign",
			input:         "KEY",
			expectedKey:   "",
			expectedValue: "",
			expectedError: true,
		},
		{
			name:          "KeyValuePairWithSpecialCharactersInValue",
			input:         "KEY=\"value_with_special#chars\"",
			expectedKey:   "KEY",
			expectedValue: "value_with_special#chars",
			expectedError: false,
		},
		{
			name:          "MultilineValue",
			input:         "KEY=\"multi\nline\nvalue\"",
			expectedKey:   "",
			expectedValue: "",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key, value, err := parseEnvFileLine(tc.input)

			// Check for unexpected errors
			if (err != nil) != tc.expectedError {
				t.Errorf("parseEnvFileLine(%q) error = %v, expectedError %v", tc.input, err, tc.expectedError)
			}

			// Verify the key
			if key != tc.expectedKey {
				t.Errorf("parseEnvFileLine(%q) key = %q, want %q", tc.input, key, tc.expectedKey)
			}

			// Verify the value
			if value != tc.expectedValue {
				t.Errorf("parseEnvFileLine(%q) value = %q, want %q", tc.input, value, tc.expectedValue)
			}
		})
	}
}

func TestReadEnvFile(t *testing.T) {
	testCases := []struct {
		name          string
		fileContent   string
		expected      map[string]string
		expectedError bool
	}{
		{
			name:        "ValidFile",
			fileContent: "KEY1=VALUE1\nKEY2=VALUE2\nKEY3=VALUE3\n",
			expected: map[string]string{
				"KEY1": "VALUE1",
				"KEY2": "VALUE2",
				"KEY3": "VALUE3",
			},
			expectedError: false,
		},
		{
			name:        "FileWithEmptyLineAndComment",
			fileContent: "# comment\n\nKEY4=VALUE4\n",
			expected: map[string]string{
				"KEY4": "VALUE4",
			},
			expectedError: false,
		},
		{
			name:          "InvalidLineWithoutEqualSign",
			fileContent:   "KEY_WITHOUT_VALUE\n",
			expected:      map[string]string{},
			expectedError: true,
		},
		{
			name:        "FileWithQoutedValue",
			fileContent: "KEY5=\"VALUE5\"",
			expected: map[string]string{
				"KEY5": "VALUE5",
			},
			expectedError: false,
		},
		{
			name:          "FileWithInvalidKey",
			fileContent:   "INVALID-KEY=VALUE\n",
			expected:      map[string]string{},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary file with the given content
			tmpFile, err := createTempFile(tc.fileContent)
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile)

			// Call the function to test
			envVars, err := readEnvFile(tmpFile)

			if (err != nil) != tc.expectedError {
				t.Errorf("readEnvFile() error = %v, expectedError %v", err, tc.expectedError)
			}

			if !tc.expectedError {
				// Verify the map content
				for key, expectedVal := range tc.expected {
					if val, exists := envVars[key]; !exists || val != expectedVal {
						t.Errorf("readEnvFile() = %v, want %s", envVars, tc.expected)
					}
				}
			}
		})
	}
}

func TestLoadEnv(t *testing.T) {
	t.Run("ValidFiles", func(t *testing.T) {
		// Create temp files with valid content
		filename1, err := createTempFile("KEY1=VALUE1\nKEY2=VALUE2\n")
		if err != nil {
			t.Fatalf("error creating temp file: %v", err)
		}
		defer os.Remove(filename1)

		filename2, err := createTempFile("KEY3=VALUE3\nKEY4=VALUE4\n")
		if err != nil {
			t.Fatalf("error creating temp file: %v", err)
		}
		defer os.Remove(filename2)

		err = loadEnv([]string{filename1, filename2})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("FileReadError", func(t *testing.T) {
		// Test non-existing file to trigger a read error
		err := loadEnv([]string{"non_existing_file.env"})
		if err == nil {
			t.Errorf("expected error loading environment file, go: %v", err)
		}
	})

	t.Run("InvalidLineError", func(t *testing.T) {
		filename, err := createTempFile("INVALID-KEY=VALUE\n")
		if err != nil {
			t.Fatalf("error creating temp file: %v", err)
		}
		defer os.Remove(filename)

		// Test invalid key
		err = loadEnv([]string{filename})
		if err == nil {
			t.Errorf("expected invalid key error, got: %v", err)
		}
	})

	t.Run("EmptyFile", func(t *testing.T) {
		// Create an empty file
		filename, err := createTempFile("")
		if err != nil {
			t.Fatalf("error creating temp file: %v", err)
		}
		defer os.Remove(filename)

		// Test with an empty file
		err = loadEnv([]string{filename})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("FileWithComments", func(t *testing.T) {
		// Create a file with comments and key-value pairs
		filename, err := createTempFile("# This is a comment\nKEY1=VALUE1\n# Another comment\nKEY2=VALUE2\n")
		if err != nil {
			t.Fatalf("error creating temp file: %v", err)
		}
		defer os.Remove(filename)

		err = loadEnv([]string{filename})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestLoad(t *testing.T) {
	// Cleanup environment variables after each test
	resetEnv := func() {
		for _, v := range []string{"API_PORT", "DB_HOST", "DB_NAME", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_SSL_MODE"} {
			os.Unsetenv(v)
		}
	}

	t.Run("ValidConfigFromEnvFile", func(t *testing.T) {
		defer resetEnv()

		// Create a temp file
		tmpFile, err := createTempFile("API_PORT=9090\nDB_HOST=localhost\nDB_NAME=testdb\nDB_PORT=5432\nDB_USER=testuser\nDB_PASSWORD=testpass\nDB_SSL_MODE=enable")
		if err != nil {
			t.Fatalf("error creating temp file: %v", err)
		}
		defer os.Remove(tmpFile)

		cfg, err := Load(tmpFile)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Validate the config
		expected := &Config{
			ApiPort:    "9090",
			DbHost:     "localhost",
			DbName:     "testdb",
			DbPort:     "5432",
			DbUser:     "testuser",
			DbPassword: "testpass",
			DbSslMode:  "enable",
		}

		if *cfg != *expected {
			t.Errorf("expected %+v, got %+v", expected, cfg)
		}
	})

	t.Run("MissingRequiredEnvVars", func(t *testing.T) {
		defer resetEnv()

		tempFile, err := createTempFile("API_PORT=8080\nDB_HOST=localhost")
		if err != nil {
			t.Fatalf("error creating temp file: %v", err)
		}
		defer os.Remove(tempFile)

		// Adjust required env variables
		requiredEnvVars = []string{"DB_NAME", "DB_USER"}

		cfg, err := Load(tempFile)
		if err == nil {
			t.Fatalf("expected error due to missing variables, got nil")
		}
		if cfg != nil {
			t.Errorf("expected nil config, got %+v", cfg)
		}
	})

	t.Run("DefaultEnvFileFallback", func(t *testing.T) {
		defer resetEnv()

		// Create the default .env file
		defaultEnvFile := ".env"
		err := os.WriteFile(defaultEnvFile, []byte("API_PORT=8000\nDB_HOST=defaultdbhost\n"), 0644)
		if err != nil {
			t.Fatalf("error creating default .env file: %v", err)
		}
		defer os.Remove(defaultEnvFile)

		requiredEnvVars = []string{"API_PORT", "DB_HOST"}

		cfg, err := Load()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Validate the config
		expected := &Config{
			ApiPort:    "8000",
			DbHost:     "defaultdbhost",
			DbName:     "",
			DbPort:     "",
			DbUser:     "",
			DbPassword: "",
			DbSslMode:  "disabled",
		}
		if *cfg != *expected {
			t.Errorf("expected %+v, got %+v", expected, cfg)
		}
	})
}

func createTempFile(content string) (string, error) {
	file, err := os.CreateTemp("", "env_file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}
