package config

type Report struct {
	ID         string      `json:"id"`
	CategoryID string      `json:"category_id"`
	Parameters []Parameter `json:"parameters"`
}

type Dataset struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Entries []struct {
	Report  Report  `json:"report"`
	Dataset Dataset `json:"dataset"`
}

type Parameter struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
