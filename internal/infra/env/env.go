package env

import (
	"os"
	"strconv"
	"time"
)

func GetEnvAsString(name string, defaultVal string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}
	return defaultVal
}

func GetEnvAsInt(name string, defaultVal int) int {
	if value, exists := os.LookupEnv(name); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultVal
}

func GetEnvAsDuration(name string, defaultVal time.Duration) time.Duration {
	if value, exists := os.LookupEnv(name); exists {
		if durationValue, err := time.ParseDuration(value); err == nil {
			return durationValue
		}
	}
	return defaultVal
}
