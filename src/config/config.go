package config

import (
	"os"
	"testing"
)

// Fetch environment variables safely
func GetEnv(key string) string {
	return os.Getenv(key)
}

func TestGetEnv(t *testing.T) {
	testKey := "TEST_ENV_VAR"
	testValue := "test_value"

	// Set the environment variable
	os.Setenv(testKey, testValue)

	// Test GetEnv function
	result := GetEnv(testKey)
	if result != testValue {
		t.Errorf("GetEnv(%s) = %s; want %s", testKey, result, testValue)
	}

	// Test non-existent key
	nonExistentKey := "NON_EXISTENT_KEY"
	result = GetEnv(nonExistentKey)
	if result != "" {
		t.Errorf("GetEnv(%s) = %s; want empty string", nonExistentKey, result)
	}
}
