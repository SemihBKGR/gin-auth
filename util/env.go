package util

import (
	"os"
	"strconv"
)

func GetEnvVar(key, def string) string {
	envVar := os.Getenv(key)
	if envVar != "" {
		return envVar
	}
	return def
}

func GetIntEnvVar(key string, def int) int {
	envVarStr := os.Getenv(key)
	if envVarStr != "" {
		envVar, err := strconv.Atoi(envVarStr)
		if err == nil {
			return envVar
		}
	}
	return def
}
