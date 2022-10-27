package config

import (
	"fmt"

	"golang.org/x/exp/slices"
)

// These map to existing servicetitan types and adds the percentage
// Other types might be added in the future as required
var validReportFieldTypes = []string{"Date", "Number", "Boolean", "String", "Percentage"}

type Report struct {
	ID         string      `yaml:"id"`
	CategoryID string      `yaml:"category_id"`
	Parameters []Parameter `yaml:"parameters"`
}

type Dataset struct {
	Name           string        `yaml:"name"`
	Type           string        `yaml:"type"`
	RequiredFields []string      `yaml:"required_fields"`
	FieldOverrides []ReportField `yaml:"field_overrides"`
}

type Entries []Entry

type Entry struct {
	Report  Report  `yaml:"report"`
	Dataset Dataset `yaml:"dataset"`
}

// ReportField allows overriding a field type of a report.
// Such as mapping a Number type to percentage
type ReportField struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type Parameter struct {
	Name  string      `yaml:"name"`
	Value interface{} `yaml:"value"`
}

func (e Entries) Validate() error {
	if len(e) == 0 {
		return Error{
			scope:    "entries",
			messages: []string{"at least one entry is required"},
		}
	}

	for idx, entry := range e {
		msgs := entry.Dataset.validate()
		msgs = append(msgs, entry.Report.validate()...)

		if len(msgs) > 0 {
			return Error{
				scope:    fmt.Sprintf("entries[%d]", idx+1),
				messages: msgs,
			}
		}
	}

	return nil
}

func (d Dataset) validate() []string {
	var msgs []string

	if len(d.RequiredFields) == 0 {
		msgs = append(msgs, "at least one dataset required_field is required, please use the report field name as the identifier")
	}

	for _, f := range d.FieldOverrides {
		if !slices.Contains(validReportFieldTypes, f.Type) {
			msg := fmt.Sprintf("field override %q type is invalid only %q are valid types", f.Name, validReportFieldTypes)
			msgs = append(msgs, msg)
		}
	}

	return msgs
}

func (r Report) validate() []string {
	var msgs []string

	if r.ID == "" {
		msgs = append(msgs, "report id is required")
	}

	if r.CategoryID == "" {
		msgs = append(msgs, "category_id is required")
	}

	return msgs
}
