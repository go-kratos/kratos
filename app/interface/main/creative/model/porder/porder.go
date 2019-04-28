package porder

// Config str
type Config struct {
	ID   int64  `json:"id"`
	Tp   int8   `json:"type"`
	Name string `json:"name"`
}

const (
	// ConfigTypeIndustry const
	ConfigTypeIndustry = 0
	// ConfigTypeShow const
	ConfigTypeShow = 1
)
