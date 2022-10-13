package config

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestError_Error(t *testing.T) {
	t.Run("returns formatted error", func(t *testing.T) {
		err := Error{
			scope:    "servicetitan",
			messages: []string{"missing attr_a", "missing attr_b"},
		}

		want := `Config section "servicetitan" errors:` + "\n - missing attr_a\n - missing attr_b"
		assert.Equal(t, err.Error(), want)
	})
}
