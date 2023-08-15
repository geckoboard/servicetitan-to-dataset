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
	var fieldType = field.Type

	// If a field override exists use that field type
	// to return the specific dataset field type
	if ovf := d.fieldOverride(field); ovf != nil {
		fieldType = ovf.Type
	}

	switch fieldType {
	case "String":
		return geckoboard.StringType
	case "Number":
		return geckoboard.NumberType
	case "Boolean":
		return geckoboard.StringType
	case "Date":
		return geckoboard.DateType
	case "Datetime":
		return geckoboard.DatetimeType
	case "Percentage":
		// Although this is a fake type override - it seems that
		// percentage fields are already a decimal from 0 to 1
		// so we don't need to do any additional conversion
		return geckoboard.PercentType
	}

	// TODO: Not sure if "Time" is a datetime or just the time ignore for now until required

	return "unknown"
}

func (d *DatasetBuilder) fieldOverride(field servicetitan.ReportField) *config.ReportField {
	for _, f := range d.datasetOverrides.FieldOverrides {
		if f.Name == field.Name {
			return &f
		}
	}

	return nil
}
