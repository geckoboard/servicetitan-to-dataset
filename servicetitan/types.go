package servicetitan

type PaginationOptions struct {
	PageSize int
	Page     int
}

type ReportDataRequest struct {
	CategoryID string               `json:"-"`
	ReportID   string               `json:"-"`
	Parameters DataRequestParamters `json:"parameters"`
}

type DataRequestParamters []struct {
	Name  string
	Value string
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CategoryList struct {
	Items []Category `json:"data"`

	HasMore  bool `json:"hasMore"`
	Page     int  `json:"page"`
	PageSize int  `json:"pageSize"`
	Total    int  `json:"totalCount"`
}

type Report struct {
	ID         int               `json:"id"`
	Name       string            `json:"name"`
	Fields     []ReportField     `json:"fields,omitempty"`
	Parameters []ReportParameter `json:"parameters,omitempty"`
}

type ReportList struct {
	Items []Report `json:"data"`

	HasMore  bool `json:"hasMore"`
	Page     int  `json:"page"`
	PageSize int  `json:"pageSize"`
	Total    int  `json:"totalCount"`
}

type ReportField struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Type  string `json:"dataType"`
}

type ReportParameter struct {
	Name       string `json:"name"`
	Label      string `json:"label"`
	DataType   string `json:"dataType"`
	IsArray    bool   `json:"isArray"`
	IsRequired bool   `json:"isRequired"`
}

type ReportData struct {
	Data   []interface{} `json:"data"`
	Fields []ReportField `json:"fields"`

	HasMore  bool `json:"hasMore"`
	Page     int  `json:"page"`
	PageSize int  `json:"pageSize"`
	Total    int  `json:"totalCount"`
}
