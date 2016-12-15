package base

// Column is a single column in data tables
type Column struct {
	Data           string            `json:"data"`
	Name           string            `json:"name"`
	Searchable     bool              `json:"searchable"`
	Sortable       bool              `json:"sortable"`
	Visible        bool              `json:"visible"`
	Filter         bool              `json:"filter"`
	Title          string            `json:"title"`
	FilterValidMap map[string]string `json:"filter_valid_map"`
}

// Columns is the columns in data tables
type Columns []Column
