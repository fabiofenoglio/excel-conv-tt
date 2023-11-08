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
	SentryDSN      string
	Environment    string
	SkipAutoUpdate bool
}

func Get() (EnvConfig, error) {
	_ = godotenv.Load(".env")

	sentryDSN := getOrDefault("SENTRY_DSN", "")

	skipAutoUpdateRaw := getOrDefault("SKIP_AUTO_UPDATE", "false")

	envName := getOrDefault("ENV", "local")

	return EnvConfig{
		SentryDSN:      sentryDSN,
		Environment:    envName,
		SkipAutoUpdate: skipAutoUpdateRaw == "true",
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
