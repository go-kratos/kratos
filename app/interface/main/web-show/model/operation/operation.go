package operation

// Types
var (
	Types = []string{
		"uptool",
		"promote",
	}
)

// Operation struct
type Operation struct {
	ID      int64  `json:"-"`
	Type    string `json:"-"`
	Message string `json:"message"`
	Link    string `json:"link"`
	Pic     string `json:"pic,omitempty"`
	Ads     int    `json:"ads,omitempty"`
	Rank    int    `json:"rank"`
	Aid     int64  `json:"-"`
}

// ArgOp ArgOp
type ArgOp struct {
	Tp    string `form:"tp"`
	Count int    `form:"count" validate:"min=0"`
	Rank  int    `form:"rank" validate:"min=0"`
}
