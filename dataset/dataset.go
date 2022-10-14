package dataset

import (
	"fmt"
	"regexp"
	"servicetitan-to-dataset/config"
	"servicetitan-to-dataset/servicetitan"
	"strconv"
	"strings"

	"github.com/jnormington/geckoboard"
)

var (
	fieldIDRegexp     = regexp.MustCompile(`[^a-z0-9 ]+|[\W]+$|^[\W]+`)
	datasetNameRegexp = regexp.MustCompile(`[^0-9a-z._\- ]+`)
)

type BuilderConfig struct {
	Report           *servicetitan.Report
	Data             *servicetitan.ReportData
	DatasetOverrides config.Dataset
}

type DatasetBuilder struct {
	report           *servicetitan.Report
	data             *servicetitan.ReportData
	datasetOverrides config.Dataset
}

func NewDatasetBuilder(conf BuilderConfig) *DatasetBuilder {
	return &DatasetBuilder{
		report:           conf.Report,
		data:             conf.Data,
		datasetOverrides: conf.DatasetOverrides,
	}
}

func (d *DatasetBuilder) BuildSchema() *geckoboard.Dataset {
	fields := map[string]geckoboard.Field{}

	uniqueFields := []string{}
	for _, f := range d.report.Fields {
		optionalField := d.isOptionalField(f)
		key := d.safeDataFieldName(f)
		fields[key] = geckoboard.Field{
			Type:     d.datasetFieldType(f),
			Name:     f.Label,
			Optional: optionalField,
		}

		if !optionalField {
			uniqueFields = append(uniqueFields, key)
		}
	}

	return &geckoboard.Dataset{
		Name:     d.datasetName(),
		Fields:   fields,
		UniqueBy: uniqueFields,
	}
}

func (d *DatasetBuilder) BuildData() geckoboard.Data {
	data := geckoboard.Data{}

	for _, r := range d.data.Data {
		gr := geckoboard.DataRow{}

		switch row := r.(type) {
		case []interface{}:
			for idx, val := range row {
				field := d.data.Fields[idx]
				name := d.safeDataFieldName(field)

				switch nval := val.(type) {
				case string:
					gr[name] = nval
				case bool:
					gr[name] = strings.ToUpper(strconv.FormatBool(nval))
				case int, float64:
					if field.Type == "String" {
						gr[name] = fmt.Sprintf("%v", nval)
					} else {
						gr[name] = nval
					}
				}
			}

			data = append(data, gr)
		default:
			panic("unexpected data row")
		}
	}

	return data
}

func (d *DatasetBuilder) safeDataFieldName(field servicetitan.ReportField) string {
	key := fieldIDRegexp.ReplaceAllString(strings.ToLower(field.Name), "")
	return strings.ReplaceAll(key, " ", "_")
}

func (d *DatasetBuilder) datasetName() string {
	name := d.report.Name

	if d.datasetOverrides.Name != "" {
		name = d.datasetOverrides.Name
	}

	nn := datasetNameRegexp.ReplaceAllLiteralString(strings.ToLower(name), "")
	return strings.ReplaceAll(nn, " ", "_")
}

func (d *DatasetBuilder) isOptionalField(field servicetitan.ReportField) bool {
	for _, rf := range d.datasetOverrides.RequiredFields {
		if rf == field.Name {
			return false
		}
	}

	return true
}

func (d *DatasetBuilder) datasetFieldType(field servicetitan.ReportField) geckoboard.FieldType {
	switch field.Type {
	case "String":
		return geckoboard.StringType
	case "Number":
		return geckoboard.NumberType
	case "Boolean":
		return geckoboard.StringType
	case "Date":
		return geckoboard.DateType
	}

	// TODO: Not sure if "Time" is a datetime or just the time ignore for now until required

	return "unknown"
}
