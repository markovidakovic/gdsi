package config

import (
	"fmt"
	"log"
	"os"

	"github.com/go-chi/jwtauth/v5"
)

type Config struct {
	ApiPort              string
	DbDriver             string
	DbHost               string
	DbName               string
	DbPort               string
	DbUser               string
	DbPassword           string
	DbSslMode            string
	JwtSecret            string
	JwtExpiration        string
	JwtAccessExpiration  string
	JwtRefreshExpiration string
	JwtAuth              *jwtauth.JWTAuth
}

const defaultEnvFile = ".env"

var requiredEnvVars = []string{"DB_DRIVER", "DB_HOST", "DB_NAME", "DB_PORT", "DB_USER", "DB_PASSWORD", "JWT_SECRET", "JWT_EXPIRATION", "JWT_ACCESS_EXPIRATION", "JWT_REFRESH_EXPIRATION"}

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
		ApiPort:              getEnvVar("API_PORT", "8080"),
		DbDriver:             getEnvVar("DB_DRIVER", ""),
		DbHost:               getEnvVar("DB_HOST", ""),
		DbName:               getEnvVar("DB_NAME", ""),
		DbPort:               getEnvVar("DB_PORT", ""),
		DbUser:               getEnvVar("DB_USER", ""),
		DbPassword:           getEnvVar("DB_PASSWORD", ""),
		DbSslMode:            getEnvVar("DB_SSL_MODE", "disabled"),
		JwtSecret:            getEnvVar("JWT_SECRET", ""),
		JwtExpiration:        getEnvVar("JWT_EXPIRATION", ""),
		JwtAccessExpiration:  getEnvVar("JWT_ACCESS_EXPIRATION", ""),
		JwtRefreshExpiration: getEnvVar("JWT_REFRESH_EXPIRATION", ""),
	}

	// add jwt auth
	cfg.JwtAuth = jwtauth.New("HS256", []byte(cfg.JwtSecret), nil)

	log.Println("config loaded")

	return cfg, nil
}
