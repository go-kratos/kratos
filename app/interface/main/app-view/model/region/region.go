package region

type Region struct {
	Rid   int    `json:"tid"`
	Reid  int    `json:"reid"`
	Name  string `json:"name"`
	Logo  string `json:"logo"`
	Goto  string `json:"goto"`
	Param string `json:"param"`
	Rank  string `json:"-"`
	Plat  int8   `json:"-"`
	Area  string `json:"-"`
}
