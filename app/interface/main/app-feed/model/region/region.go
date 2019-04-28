package region

type Region struct {
	Rid       int64     `json:"tid"`
	Reid      int64     `json:"reid"`
	Name      string    `json:"name"`
	Logo      string    `json:"logo"`
	Goto      string    `json:"goto"`
	Param     string    `json:"param"`
	Rank      string    `json:"-"`
	Plat      int8      `json:"-"`
	Build     int       `json:"-"`
	Condition string    `json:"-"`
	Area      string    `json:"-"`
	Language  string    `json:"-"`
	Children  []*Region `json:"children,omitempty"`
}
