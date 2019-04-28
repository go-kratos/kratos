package app

import (
	"go-common/library/time"
)

// platMap and Article vars
var (
	platMap = map[string][]int8{
		"android": {0, 1},
		"ios":     {0, 2},
		"ipad":    {0, 3},
	}
	MyArticle    = "我的专栏"
	OpenArticle  = "开通专栏"
	PortalIntro  = 0
	PortalNotice = 1
	BgmFrom      = map[int]int{
		0: 1, // videoup
		1: 1, // camera
	}
)

// Type
const (
	TypeSwitch         = int8(-1)
	TypeSubtitle       = int8(0)
	TypeFont           = int8(1)
	TypeFilter         = int8(2)
	TypeBGM            = int8(3)
	TypeHotWord        = int8(4)
	TypeSticker        = int8(5)
	TypeIntro          = int8(6)
	TypeVideoupSticker = int8(7)
	TypeTransition     = int8(8)
	TypeCooperate      = int8(9)  // 拍摄之稿件合拍
	TypeTheme          = int8(10) // 编辑器的主题使用相关

	//For portal White category
	WhiteGroupType   = int8(1)
	WhitePercentType = int8(2)
	WhitePercentV00  = int64(0)
	WhitePercentV10  = int64(1)
	WhitePercentV20  = int64(2)
	WhitePercentV50  = int64(5)
)

// Icon str
type Icon struct {
	State bool   `json:"state"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

// AllowMaterial fn
func (pm *PortalMeta) AllowMaterial(platStr string, build int) (ret bool) {
	if pm.Platform == 0 {
		return true
	}
	platOK := false
	for _, num := range platMap[platStr] {
		if pm.Platform == num {
			platOK = true
		}
	}
	buildOK := true
	for _, v := range pm.BuildExps {
		if !AllowBuild(build, v.Condition, v.Build) {
			buildOK = false
		}
	}
	return buildOK && platOK
}

// AcademyIntro str
type AcademyIntro struct {
	Title string `json:"title"`
	Show  int    `json:"show"`
	URL   string `json:"url"`
}

// ActIntro str
type ActIntro struct {
	Title string `json:"title"`
	Show  int    `json:"show"`
	URL   string `json:"url"`
}

//Portal for app portal
type Portal struct {
	ID       int64     `json:"id"`
	Icon     string    `json:"icon"`
	Title    string    `json:"title"`
	Pos      int8      `json:"position"`
	URL      string    `json:"url"`
	New      int8      `json:"new"`
	More     int8      `json:"more"`
	SubTitle string    `json:"subtitle"`
	MTime    time.Time `json:"mtime"`
}

//Up verify author.
type Up struct {
	Arc int `json:"archive"`
	Art int `json:"article"`
}

// BuildExp str
type BuildExp struct {
	Condition int8 `json:"conditions"`
	Build     int  `json:"build"`
}

// WhiteExp str
type WhiteExp struct {
	TP    int8  `json:"type"`
	Value int64 `json:"value"`
}

//PortalMeta for app.
type PortalMeta struct {
	ID        int64       `json:"id"`
	Build     int         `json:"build"`
	Platform  int8        `json:"platform"` ////0-全平台,1-android,2-ios,3-ipad,4-ipod
	Compare   int8        `json:"compare"`  //比较版本号符号类型,0-等于,1-小于,2-大于,3-不等于,4-小于等于,5-大于等于
	State     int8        `json:"state"`    //状态,0-关闭,1-打开
	Pos       int8        `json:"pos"`
	Mark      int8        `json:"mark"` //是否为新的标记,0-否,1-是
	Type      int8        `json:"type"` //业务类型,0-创作中心,1-主站APP个人中心配置
	Title     string      `json:"title"`
	Icon      string      `json:"icon"`
	URL       string      `json:"url"`
	CTime     time.Time   `json:"ctime"`
	MTime     time.Time   `json:"mtime"`
	PTime     time.Time   `json:"ptime"`
	More      int8        `json:"more"`
	BuildExps []*BuildExp `json:"-"`
	// add from app 537
	SubTitle  string      `json:"subtitle"`
	WhiteExps []*WhiteExp `json:"-"`
}

//PlatFormMap for  platform from int to string map.
func PlatFormMap(pf int8) (res string) {
	switch pf {
	case 1:
		res = "android"
	case 2:
		res = "ios"
	}
	return
}

//EarningsCopyWriter for creator earnings copywriter.
type EarningsCopyWriter struct {
	Elec   string `json:"elec"`
	Growth string `json:"growth"`
	Oasis  string `json:"oasis"`
}

//BuildMap for build from int to string map.
func BuildMap(pf int8) (res string) { //比较版本号符号类型,0-等于,1-小于,2-大于,3-不等于,4-小于等于,5-大于等于
	switch pf {
	case 0:
		res = "="
	case 1:
		res = "<"
	case 2:
		res = ">"
	case 3:
		res = "!="
	case 4:
		res = "<="
	case 5:
		res = ">="
	}
	return
}

// AllowPlatForm fn
func AllowPlatForm(platStr string, plat int8) bool {
	if plat == 0 || platStr == PlatFormMap(plat) {
		return true
	}
	return false
}

// AllowBuild fn
func AllowBuild(buildParam int, compare int8, build int) (res bool) {
	switch BuildMap(compare) {
	case "=":
		if buildParam == build {
			res = true
		}
	case "<":
		if buildParam < build {
			res = true
		}
	case ">":
		if buildParam > build {
			res = true
		}
	case "!=":
		if buildParam != build {
			res = true
		}
	case "<=":
		if buildParam <= build {
			res = true
		}
	case ">=":
		if buildParam >= build {
			res = true
		}
	}
	return
}
