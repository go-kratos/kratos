package music

import (
	"encoding/json"
	"go-common/app/interface/main/creative/model/activity"
	"go-common/app/interface/main/creative/model/app"
	accMdl "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	"go-common/library/time"
)

var (
	platMap = map[string][]int{
		"android": {0, 1},
		"ios":     {0, 2},
	}
	// ViewTpMap map
	ViewTpMap = map[int8]string{
		0:  "subtitle",
		1:  "font",
		2:  "filter",
		5:  "sticker",
		7:  "videoup_sticker",
		8:  "transition",
		9:  "cooperate",
		10: "theme",
	}
)

// BuildComp str
type BuildComp struct {
	Condition int8 `json:"conditions"`
	Build     int  `json:"build"`
}

// AllowMaterial fn
// 1:platform first; 2:build alg; 3:user whitelist
func (v *Material) AllowMaterial(m Material, platStr string, buildParam int, white bool) (ret bool) {
	if v.White == 1 && !white {
		return false
	}
	if v.Platform == 0 {
		return true
	}
	platOK := false
	for _, num := range platMap[platStr] {
		if m.Platform == num {
			platOK = true
		}
	}
	buildOK := true
	for _, v := range m.BuildComps {
		if !app.AllowBuild(buildParam, v.Condition, v.Build) {
			buildOK = false
		}
	}
	return buildOK && platOK
}

// Music str
type Music struct {
	ID             int64           `json:"id"`
	TID            int             `json:"tid"`
	Index          int             `json:"index"`
	SID            int64           `json:"sid"`
	Name           string          `json:"name"`
	Musicians      string          `json:"musicians"`
	UpMID          int64           `json:"mid"`
	Cover          string          `json:"cover"`
	Stat           string          `json:"stat"`
	Playurl        string          `json:"playurl"`
	State          int             `json:"state"`
	Duration       int             `json:"duration"`
	FileSize       int             `json:"filesize"`
	CTime          time.Time       `json:"ctime"`
	Pubtime        time.Time       `json:"pubtime"`
	MTime          time.Time       `json:"mtime"`
	TagsStr        string          `json:"-"`
	Tags           []string        `json:"tags"`
	Timeline       json.RawMessage `json:"-"`
	Tl             []*TimePoint    `json:"timeline"`
	RecommendPoint int64           `json:"recommend_point"`
	Cooperate      int8            `json:"cooperate"`
	CooperateURL   string          `json:"cooperate_url"`
	New            int8            `json:"new"`
	Hotval         int             `json:"hotval"`
}

// BgmExt str
type BgmExt struct {
	Msc          *Music          `json:"msc"`
	ExtMscs      []*Music        `json:"ext_mscs"`
	ExtArcs      []*api.Arc      `json:"ext_arcs"`
	UpProfile    *accMdl.Profile `json:"up_profile"`
	ShouldFollow bool            `json:"show_follow"`
}

// TimePoint str
type TimePoint struct {
	Point     int64  `json:"point"`
	Comment   string `json:"comment"`
	Recommend int    `json:"recommend"`
}

// Category str
type Category struct {
	ID          int      `json:"id"`
	PID         int      `json:"pid"`
	Name        string   `json:"name"`
	Index       int      `json:"index"`
	CameraIndex int      `json:"camera_index"`
	Children    []*Music `json:"children"`
}

// Mcategory str
type Mcategory struct {
	ID    int64     `json:"id"`
	SID   int64     `json:"sid"`
	Tid   int       `json:"tid"`
	Index int       `json:"index"`
	CTime time.Time `json:"ctime"`
	New   int8      `json:"new"`
}

// Audio str
type Audio struct {
	Title string `json:"title"`
	Cover string `json:"cover_url"`
}

// Material str
type Material struct {
	Type       int8            `json:"type"`
	Platform   int             `json:"platform"`
	Build      json.RawMessage `json:"build"`
	BuildComps []*BuildComp    `json:"build_comps"`
	White      int8            `json:"white"`
	New        int8            `json:"new"`
}

// Basic str
type Basic struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Cover       string          `json:"cover"`
	DownloadURL string          `json:"download_url"`
	Rank        int             `json:"rank"`
	Max         int             `json:"max"`
	Extra       json.RawMessage `json:"-"`
	Material    `json:"-"`
	MTime       time.Time `json:"mtime"`
	Tags        []string  `json:"tags"`
}

// Cooperate str db+search+api
type Cooperate struct {
	ID       int64           `json:"id"`
	Name     string          `json:"name"`
	Cover    string          `json:"cover"`
	Rank     int             `json:"rank"`
	Extra    json.RawMessage `json:"-"`
	Material `json:"-"`
	MTime    time.Time `json:"mtime"`
	New      int8      `json:"new"`
	Tags     []string  `json:"tags"`
	// special extra column for cooperate
	MaterialAID int64              `json:"material_aid"`
	MaterialCID int64              `json:"material_cid"`
	DemoAID     int64              `json:"demo_aid"`
	DemoCID     int64              `json:"demo_cid"`
	MissionID   int64              `json:"mission_id"`
	SubType     int                `json:"sub_type"`
	Style       int                `json:"style"`
	Mission     *activity.Activity `json:"mission_info"`
	HotVal      int                `json:"hotval"`
	ArcCnt      int                `json:"-"`
	DownloadURL string             `json:"download_url"`
}

// Subtitle str
type Subtitle struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Cover       string          `json:"cover"`
	DownloadURL string          `json:"download_url"`
	Rank        int             `json:"rank"`
	Max         int             `json:"max"`
	Extra       json.RawMessage `json:"-"`
	Material    `json:"-"`
	MTime       time.Time `json:"mtime"`
	New         int8      `json:"new"`
	Tags        []string  `json:"tags"`
}

// Font str
type Font struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Cover       string          `json:"cover"`
	DownloadURL string          `json:"download_url"`
	Rank        int             `json:"rank"`
	Extra       json.RawMessage `json:"-"`
	Material    `json:"-"`
	MTime       time.Time `json:"mtime"`
	New         int8      `json:"new"`
	Tags        []string  `json:"tags"`
}

// Filter str
type Filter struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Cover       string          `json:"cover"`
	DownloadURL string          `json:"download_url"`
	Rank        int             `json:"rank"`
	Extra       json.RawMessage `json:"-"`
	Material    `json:"-"`
	MTime       time.Time `json:"mtime"`
	New         int8      `json:"new"`
	Tags        []string  `json:"tags"`
	FilterType  int8      `json:"filter_type"`
}

// VSticker str
type VSticker struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Cover       string          `json:"cover"`
	DownloadURL string          `json:"download_url"`
	Rank        int             `json:"rank"`
	Extra       json.RawMessage `json:"-"`
	Material    `json:"-"`
	MTime       time.Time `json:"mtime"`
	New         int8      `json:"new"`
	Tags        []string  `json:"tags"`
}

// Transition str
type Transition struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Cover       string          `json:"cover"`
	DownloadURL string          `json:"download_url"`
	Rank        int             `json:"rank"`
	Extra       json.RawMessage `json:"-"`
	Material    `json:"-"`
	MTime       time.Time `json:"mtime"`
	New         int8      `json:"new"`
	Tags        []string  `json:"tags"`
}

// Theme str
type Theme struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Cover       string          `json:"cover"`
	DownloadURL string          `json:"download_url"`
	Rank        int             `json:"rank"`
	Extra       json.RawMessage `json:"-"`
	Material    `json:"-"`
	MTime       time.Time `json:"mtime"`
	New         int8      `json:"new"`
	Tags        []string  `json:"tags"`
}

// Sticker str
type Sticker struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Cover       string          `json:"cover"`
	DownloadURL string          `json:"download_url"`
	Rank        int             `json:"rank"`
	Extra       json.RawMessage `json:"-"`
	Material    `json:"-"`
	SubType     int64     `json:"sub_type"`
	Tip         string    `json:"tip"`
	MTime       time.Time `json:"mtime"`
	New         int8      `json:"new"`
	Tags        []string  `json:"tags"`
}

// Hotword str
type Hotword struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Cover       string          `json:"cover"`
	DownloadURL string          `json:"download_url"`
	Rank        int             `json:"rank"`
	Extra       json.RawMessage `json:"-"`
	Material    `json:"-"`
	MTime       time.Time `json:"mtime"`
	New         int8      `json:"new"`
	Tags        []string  `json:"tags"`
}

// Intro str
type Intro struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Cover       string          `json:"cover"`
	DownloadURL string          `json:"download_url"`
	Rank        int             `json:"rank"`
	Extra       json.RawMessage `json:"-"`
	Material    `json:"-"`
	MTime       time.Time `json:"mtime"`
	New         int8      `json:"new"`
	Tags        []string  `json:"tags"`
}

// MaterialBind str
type MaterialBind struct {
	CID   int64
	MID   int64
	CName string
	CRank int
	BRank int
	Tp    int
	New   int
}

// FilterCategory str
type FilterCategory struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Rank     int       `json:"rank"`
	Tp       int       `json:"type"`
	Children []*Filter `json:"children"`
	New      int       `json:"new"`
}

// VstickerCategory str
type VstickerCategory struct {
	ID       int64       `json:"id"`
	Name     string      `json:"name"`
	Rank     int         `json:"rank"`
	Tp       int         `json:"type"`
	Children []*VSticker `json:"children"`
	New      int         `json:"new"`
}
