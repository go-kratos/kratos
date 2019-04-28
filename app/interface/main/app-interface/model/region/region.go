package region

type Region struct {
	Rid       int16  `json:"tid"`
	Reid      int16  `json:"reid"`
	Name      string `json:"name"`
	Logo      string `json:"logo"`
	Goto      string `json:"goto"`
	Param     string `json:"param"`
	Rank      string `json:"-"`
	Plat      int8   `json:"-"`
	Build     int    `json:"-"`
	Condition string `json:"-"`
	Area      string `json:"-"`
}
