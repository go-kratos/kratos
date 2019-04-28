package app

var (
	// PlatformMap map
	PlatformMap = map[int8]string{
		0: "全平台",
		1: "Android",
		2: "iOS",
		3: "iPad",
	}
)

//Portal for app.
type Portal struct {
	ID       int64  `form:"id" json:"id" gorm:"primary_key"`
	Build    int64  `form:"build" json:"build"`
	BuildExp string `form:"buildexp" json:"buildexp" gorm:"column:buildexp"`
	Platform int8   `form:"platform" json:"platform" gorm:"column:platform"`
	Compare  int8   `form:"compare" json:"compare"`
	State    int8   `form:"state" json:"state" gorm:"column:state"`
	Pos      int16  `form:"pos" json:"pos"`
	Mark     int8   `form:"mark" json:"mark"`
	More     int8   `form:"more" json:"more"`
	Type     int8   `form:"type" json:"type" gorm:"column:type"`
	Title    string `form:"title" json:"title"`
	Icon     string `form:"icon" json:"icon"`
	URL      string `form:"url" json:"url"`
	CTime    string `form:"ctime" json:"ctime" gorm:"column:ctime"`
	MTime    string `form:"mtime" json:"mtime" gorm:"column:mtime"`
	PTime    string `form:"ptime" json:"ptime" gorm:"column:ptime"`
	SubTitle string `form:"subtitle" json:"subtitle" gorm:"column:subtitle"`
	WhiteExp string `form:"whiteexp" json:"whiteexp" gorm:"column:whiteexp"`
}

// PortalPager def.
type PortalPager struct {
	Total int64   `json:"total"`
	Pn    int     `json:"pn"`
	Ps    int     `json:"ps"`
	Items []*Item `json:"items"`
}

// WhiteExp str
type WhiteExp struct {
	TP    int8 `json:"type"`
	Value int  `json:"value"`
}

//Item for portal list.
type Item struct {
	ID        int64       `json:"id"`
	Build     int64       `json:"build"`
	BuildExp  string      `json:"buildexp"`
	Platform  int8        `json:"platform"`
	Compare   int8        `json:"compare"`
	State     int8        `json:"state"`
	Pos       int16       `json:"pos"`
	Mark      int8        `json:"mark"`
	More      int8        `json:"more"`
	Type      int8        `json:"type"`
	Title     string      `json:"title"`
	Icon      string      `json:"icon"`
	URL       string      `json:"url"`
	CTime     int64       `json:"ctime"`
	MTime     int64       `json:"mtime"`
	PTime     int64       `json:"ptime"`
	SubTitle  string      `json:"subtitle"`
	WhiteExps []*WhiteExp `json:"whiteexp"`
}

// TableName fn
func (Portal) TableName() string {
	return "app_portal"
}
