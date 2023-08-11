package config

import (
	"fmt"
	"os"
	"runtime"
)

func GetSystemEnvVar(envVar string) (string, error) {
	// Retrieve the operating system
	ros := runtime.GOOS

	// Retrieve the environment variable value based on the operating system
	var value string
	switch ros {
	case "windows":
		// For Windows, use %VAR_NAME% syntax
		value = os.Getenv(envVar)
	case "linux", "darwin":
		// For Linux and macOS, use $VAR_NAME syntax
		value = os.Getenv(envVar)
	default:
		return "", fmt.Errorf("unsupported operating system: %s", ros)
	}

	// Check if the environment variable is empty
	if value == "" {
		return "", fmt.Errorf("environment variable not found: %s", envVar)
	}

	return value, nil
}
