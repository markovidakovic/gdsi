// config_test.go contains unit tests for the config package.
// These tests validate the functionality for managing application
// configuration, including loading envrionment variables from files,
// retrieving configuration values with defaults, and handling errors
// for missing or invalid environment variables
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
