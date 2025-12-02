package env

import (
	"greddit/internal/infra/log"
	"os"
)

var (
	// Over here, we use the default logger, but it can be swapped out for another safer one, if the default
	// logger is overridden with one which depends on an environment variable.
	logger = log.GetDefaultLogger()
)

// GetEnvOrDef returns the environment variable with the given key or the
// default value if the environment variable is not set.
func GetEnvOrDef[T any](key string, defV T, adapters ...Adapter[T]) T {
	str, ok := os.LookupEnv(key)
	if !ok {
		return defV
	}

	for _, adapter := range adapters {
		v, err := adapter(str)
		if err == nil {
			return v
		}
	}

	logger.Warn("Unable to adapt environment variable, using default value",
		"key", key, "default", defV)
	return defV
}

// GetEnvOrFatal returns the environment variable with the given key or exits
// the program if the environment variable is not set.
func GetEnvOrFatal[T any](key string, adapters ...Adapter[T]) T {
	str, ok := os.LookupEnv(key)
	if !ok {
		logger.Error("Required environment variable not set, exiting", "key", key)
		os.Exit(1)
	}
	for _, adapter := range adapters {
		v, err := adapter(str)
		if err == nil {
			return v
		}
	}

	logger.Error("Required environment variable not set, exiting", "key", key)
	os.Exit(1)
	panic("Should not reach here")
}

// GetStringEnvOrFatal is GetEnvOrFatal for string values.
func GetStringEnvOrFatal(key string) string {
	return GetEnvOrFatal(key, StrAdapter)
}

// GetStringEnvDef is GetEnvOrDef for string values.
func GetStringEnvDef(key string, defV string) string {
	return GetEnvOrDef(key, defV, StrAdapter)
}

// GetIntEnvOrFatal is GetEnvOrFatal for int values.
func GetIntEnvOrFatal(key string) int {
	return GetEnvOrFatal(key, IntAdapter)
}

// GetIntEnvDef is GetEnvOrDef for int values.
func GetIntEnvDef(key string, defV int) int {
	return GetEnvOrDef(key, defV, IntAdapter)
}

// GetF64EnvOrFatal is GetEnvOrFatal for float64 values.
func GetF64EnvOrFatal(key string) float64 {
	return GetEnvOrFatal(key, F64Adapter)
}

// GetF64EnvDef is GetEnvOrDef for float64 values.
func GetF64EnvDef(key string, defV float64) float64 {
	return GetEnvOrDef(key, defV, F64Adapter)
}

// GetBoolEnvOrFatal is GetEnvOrFatal for bool values.
func GetBoolEnvOrFatal(key string) bool {
	return GetEnvOrFatal(key, BoolAdapter)
}

// GetBoolEnvDef is GetEnvOrDef for bool values.
func GetBoolEnvDef(key string, defV bool) bool {
	return GetEnvOrDef(key, defV, BoolAdapter)
}
