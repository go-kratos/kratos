package pgc

// ParamStyle .
type ParamStyle struct {
	Name    string `json:"name"`
	StyleID int    `json:"style_id"`
}

// StyleRes .
type StyleRes struct {
	ID       int
	Style    string
	Category int
}

// LabelRes .
type LabelRes struct {
	Name     string
	Value    int
	Category int
}
