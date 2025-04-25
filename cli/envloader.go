package sdk

import (
	"os"
)

// Env retrieves environment variables safely.
func Env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
