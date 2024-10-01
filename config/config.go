package config

import (
	"os"
)

// Fetch environment variables safely
func GetEnv(key string) string {
	return os.Getenv(key)
}
