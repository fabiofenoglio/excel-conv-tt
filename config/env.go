package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
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
