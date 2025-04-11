// Package env provides functions for retrieving values from environment variables with fallback values.
package env

import (
	"os"
	"strconv"
	"time"
)

// GetString retrives the value of environment variable `k`. If no value is found, then the fallback value is returned.
func GetString(k, f string) string {
	v, found := os.LookupEnv(k)
	if !found {
		return f
	}
	return v
}

// GetDuration retrieves the value of environment variable `k`. If the variable is present, then the value is parsed
// into `time.Duration`. Should this fail, then the fallback value is returned. If the variable is not present, then
// the fallback value is returned.
func GetDuration(k string, f time.Duration) time.Duration {
	v, found := os.LookupEnv(k)
	if !found {
		return f
	}
	dur, err := time.ParseDuration(v)
	if err != nil {
		return f
	}
	return dur
}

// GetInt retrieves the value of environment variable `k`. If the variable is present, then the value is parsed
// into type `int`. If the variable is not present, then the fallback is returned.
func GetInt(k string, f int) int {
	v, found := os.LookupEnv(k)
	if !found {
		return f
	}
	int, err := strconv.Atoi(v)
	if err != nil {
		return f
	}
	return int
}
