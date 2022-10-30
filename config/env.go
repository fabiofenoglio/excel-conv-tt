package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

var (
	injected = map[string]string{
		"SENTRY_DSN": "___SENTRY_DSN___",
		"ENV":        "___ENV___",
	}
)

type EnvConfig struct {
	SentryDSN   string
	Environment string
}

func Get() (EnvConfig, error) {
	_ = godotenv.Load(".env")

	sentryDSN, err := get("SENTRY_DSN")
	if err != nil {
		return EnvConfig{}, err
	}

	envName, err := get("ENV")
	if err != nil {
		return EnvConfig{}, err
	}

	return EnvConfig{
		SentryDSN:   sentryDSN,
		Environment: envName,
	}, nil
}

func get(key string) (string, error) {
	v := getOrDefault(key, "")
	if v == "" {
		return "", errors.New("missing configuration value: " + key)
	}
	return v, nil
}

func getOrDefault(key string, defaultValue string) string {
	injectedValue := injected[key]
	if injectedValue == "" || strings.HasPrefix(injectedValue, "___") {
		injectedValue = os.Getenv(key)
	}
	if injectedValue == "" {
		return defaultValue
	}
	return injectedValue
}
