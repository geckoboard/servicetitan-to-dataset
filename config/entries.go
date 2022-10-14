package config

import "fmt"

type Report struct {
	ID         string      `yaml:"id"`
	CategoryID string      `yaml:"category_id"`
	Parameters []Parameter `yaml:"parameters"`
}

type Dataset struct {
	Name           string   `yaml:"name"`
	Type           string   `yaml:"type"`
	RequiredFields []string `yaml:"required_fields"`
}

type Entries []Entry

type Entry struct {
	Report  Report  `yaml:"report"`
	Dataset Dataset `yaml:"dataset"`
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
