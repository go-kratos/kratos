package jpush

// Audience .
type Audience struct {
	Tag     interface{} `json:"tag,omitempty"`
	TagAnd  interface{} `json:"tag_and,omitempty"`
	TagNot  interface{} `json:"tag_not,omitempty"`
	Alias   interface{} `json:"alias,omitempty"`
	RegID   interface{} `json:"registration_id,omitempty"`
	Segment interface{} `json:"segment,omitempty"`
	AbTest  interface{} `json:"abtest,omitempty"`
}
