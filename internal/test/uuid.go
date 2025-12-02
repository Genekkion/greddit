package test

import (
	"github.com/google/uuid"
)

// StringsToUuids converts a slice of strings to a slice of UUIDs. Only for
// testing purposes.
func StringsToUuids(arr []string) []uuid.UUID {
	res := make([]uuid.UUID, 0, len(arr))

	for _, id := range arr {
		res = append(res,
			uuid.MustParse(id),
		)
	}

	return res
}
