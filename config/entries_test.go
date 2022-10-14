package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert"
)

func TestEntries_Validate(t *testing.T) {
	t.Run("returns error at least one entry required", func(t *testing.T) {
		want := Error{
			scope:    "entries",
			messages: []string{"at least one entry is required"},
		}

		in := Entries{}
		assert.DeepEqual(t, in.Validate(), want, cmp.AllowUnexported(Error{}))
	})

	t.Run("returns all errors from dataset and report", func(t *testing.T) {
		want := Error{
			scope: "entries[1]",
			messages: []string{
				"at least one dataset required_field is required, please use the report field name as the identifier",
				"report id is required",
				"category_id is required",
			},
		}

		in := Entries{{}}
		assert.DeepEqual(t, in.Validate(), want, cmp.AllowUnexported(Error{}))
	})

	t.Run("returns specific errors report", func(t *testing.T) {
		want := Error{
			scope: "entries[1]",
			messages: []string{
				"report id is required",
				"category_id is required",
			},
		}

		in := Entries{{Dataset: Dataset{RequiredFields: []string{"Name"}}}}
		assert.DeepEqual(t, in.Validate(), want, cmp.AllowUnexported(Error{}))
	})

	t.Run("returns error for later entry", func(t *testing.T) {
		want := Error{
			scope: "entries[3]",
			messages: []string{
				"at least one dataset required_field is required, please use the report field name as the identifier",
			},
		}

		in := Entries{
			{
				Dataset: Dataset{
					RequiredFields: []string{"Name"},
				},
				Report: Report{
					ID:         "rpt-1",
					CategoryID: "cat-1",
				},
			},
			{
				Dataset: Dataset{
					RequiredFields: []string{"Name"},
				},
				Report: Report{
					ID:         "rpt-3",
					CategoryID: "cat-3",
				},
			},
			{
				Dataset: Dataset{},
				Report: Report{
					ID:         "rpt-5",
					CategoryID: "cat-5",
				},
			},
		}
		assert.DeepEqual(t, in.Validate(), want, cmp.AllowUnexported(Error{}))
	})

	t.Run("returns no errors", func(t *testing.T) {
		in := Entries{
			{
				Dataset: Dataset{
					RequiredFields: []string{"Name"},
				},
				Report: Report{
					ID:         "rpt-1",
					CategoryID: "cat-1",
				},
			},
			{
				Dataset: Dataset{
					RequiredFields: []string{"Name"},
				},
				Report: Report{
					ID:         "rpt-3",
					CategoryID: "cat-3",
				},
			},
			{
				Dataset: Dataset{
					RequiredFields: []string{"Name"},
				},
				Report: Report{
					ID:         "rpt-5",
					CategoryID: "cat-5",
				},
			},
		}

		assert.NilError(t, in.Validate())
	})
}
