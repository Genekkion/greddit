package test

import (
	"encoding/json"

	"greddit/internal/util"
)

// JsonB is a simple wrapper around json.Marshal that panics on error.
func JsonB(v any) []byte {
	return util.Must(func() ([]byte, error) {
		return json.Marshal(v)
	})
}

// JsonS is a simple wrapper around json.Marshal that panics on error and returns a string.
func JsonS(v any) string {
	return string(JsonB(v))
}
