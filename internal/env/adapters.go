package env

import (
	"strconv"
	"strings"
)

/*
	The adapters are used to parse environment variables into the desired type. The
	adapters are mainly used in the GetEnv functions.
*/

// Adapter is a function that adapts a string to a type.
type Adapter[T any] func(str string) (T, error)

// IntAdapter parses an integer from a string.
func IntAdapter(str string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(str))
}

// F64Adapter parses a float64 from a string.
func F64Adapter(str string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(str), 64)
}

// StrAdapter returns the string as is.
func StrAdapter(str string) (string, error) {
	return str, nil
}

// BoolAdapter parses a boolean from a string.
func BoolAdapter(str string) (bool, error) {
	return strconv.ParseBool(str)
}
