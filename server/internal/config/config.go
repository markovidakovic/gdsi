package config

import (
	"fmt"
	"log"
	"os"

	"github.com/go-chi/jwtauth/v5"
)

// Config holds application-specific configuration values,
// such as API server port and database connection details.
type Config struct {
	ApiPort       string           // Port where the API server runs
	DbDriver      string           // Database driver
	DbHost        string           // Host of the database
	DbName        string           // Name of the database
	DbPort        string           // Port of the database connection
	DbUser        string           // Database username
	DbPassword    string           // Database password
	DbSslMode     string           // SSL mode for database connection
	JwtSecret     string           // Json web token secret
	JwtExpiration string           // Json web token expiration
	JwtAuth       *jwtauth.JWTAuth // go-chi/jwtauth/v5 jwt auth type
}

// Default file path for loading environment variables.
const defaultEnvFile = ".env"

// requiredEnvVars lists the environment variables that must be set for the
// application to run correctly. An error is returned if any of these are missing.
var requiredEnvVars = []string{"DB_DRIVER", "DB_HOST", "DB_NAME", "DB_PORT", "DB_USER", "DB_PASSWORD", "JWT_SECRET", "JWT_EXPIRATION"}

// Load reads environment variables from specified env files or defaults to the
// ".env" file. It returns a pointer to a config struct populated with values or an error
// if required variables are missing
func Load(envFiles ...string) (*Config, error) {
	if len(envFiles) == 0 {
		envFiles = append(envFiles, defaultEnvFile)
	}

	err := loadEnv(envFiles)
	if err != nil {
		return nil, err
	}

	for _, rev := range requiredEnvVars {
		if os.Getenv(rev) == "" {
			return nil, fmt.Errorf("missing environment variable: %s", rev)
		}
	}

	var cfg *Config = &Config{
		ApiPort:       getEnvVar("API_PORT", "8080"),
		DbDriver:      getEnvVar("DB_DRIVER", ""),
		DbHost:        getEnvVar("DB_HOST", ""),
		DbName:        getEnvVar("DB_NAME", ""),
		DbPort:        getEnvVar("DB_PORT", ""),
		DbUser:        getEnvVar("DB_USER", ""),
		DbPassword:    getEnvVar("DB_PASSWORD", ""),
		DbSslMode:     getEnvVar("DB_SSL_MODE", "disabled"),
		JwtSecret:     getEnvVar("JWT_SECRET", ""),
		JwtExpiration: getEnvVar("JWT_EXPIRATION", ""),
	}

	// Add jwt auth
	cfg.JwtAuth = jwtauth.New("HS256", []byte(cfg.JwtSecret), nil)

	log.Println("config loaded")

	return cfg, nil
}
