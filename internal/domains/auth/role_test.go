package auth

import (
	"greddit/internal/test"
	"testing"
)

func TestRoleValidate(t *testing.T) {
	t.Run("valid roles", func(t *testing.T) {
		t.Parallel()

		data := []string{
			"admin",
			"user",
		}
		for _, v := range data {
			role := Role(v)
			test.NilErr(t, role.Validate())
		}
	})

	t.Run("invalid roles", func(t *testing.T) {
		t.Parallel()

		data := []string{
			"admin1",
			"user1",
			"",
		}
		for _, v := range data {
			role := Role(v)
			err := role.Validate()
			test.Assert(t, "Validate should flag an error", err != nil)
		}
	})
}
