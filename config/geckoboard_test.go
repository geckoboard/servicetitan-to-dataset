package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert"
)

func TestGeckoboard_Validate(t *testing.T) {
	t.Run("returns error missing api_key", func(t *testing.T) {
		want := Error{
			scope:    "geckoboard",
			messages: []string{"missing api_key"},
		}

		in := Geckoboard{}
		assert.DeepEqual(t, in.Validate(), want, cmp.AllowUnexported(Error{}))
	})

	t.Run("returns no error when all valid", func(t *testing.T) {
		in := Geckoboard{APIKey: "api123"}
		assert.NilError(t, in.Validate())
	})
}
