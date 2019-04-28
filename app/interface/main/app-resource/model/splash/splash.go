package splash

import (
	"encoding/json"
	"go-common/app/interface/main/app-resource/model"
	xtime "go-common/library/time"
)

// Splash is splash type.
type Splash struct {
	ID        int64      `json:"id"`
	Type      int8       `json:"type"`
	Animate   int8       `json:"animate"`
	Duration  int16      `json:"duration"`
	Start     xtime.Time `json:"start_time,omitempty"`
	End       xtime.Time `json:"end_time,omitempty"`
	Image     string     `json:"thumb"`
	Hash      string     `json:"hash"`
	Times     int16      `json:"times"`
	Skip      int8       `json:"skip"`
	URI       string     `json:"uri"`
	Area      string     `json:"-"`
	Plat      int8       `json:"-"`
	Goto      string     `json:"-"`
	Param     string     `json:"-"`
	Width     int        `json:"-"`
	Height    int        `json:"-"`
	Build     int        `json:"-"`
	Condition string     `json:"-"`
	Operate   int        `json:"-"`
	NoPreview int        `json:"-"`
	// bitrhday
	BirthStart      string `json:"start_date,omitempty"`
	BirthEnd        string `json:"end_date,omitempty"`
	BirthStartMonth string `json:"-"`
	BirthEndMonth   string `json:"-"`
}

type List struct {
	ID                    int64           `json:"id"`
	Type                  int8            `json:"type"`
	CardType              int8            `json:"card_type"`
	Duration              int16           `json:"duration"`
	Start                 xtime.Time      `json:"begin_time,omitempty"`
	End                   xtime.Time      `json:"end_time,omitempty"`
	Image                 string          `json:"thumb"`
	Hash                  string          `json:"hash"`
	LogoURL               string          `json:"logo_url"`
	LogoHash              string          `json:"logo_hash"`
	Skip                  int8            `json:"skip"`
	URI                   string          `json:"uri"`
	VideoURL              string          `json:"video_url,omitempty"`
	VideoHash             string          `json:"video_hash,omitempty"`
	VideoWidth            int             `json:"video_width,omitempty"`
	VideoHeight           int             `json:"video_height,omitempty"`
	URITitle              string          `json:"uri_title"`
	Source                int             `json:"source,omitempty"`
	CmMark                int             `json:"cm_mark,omitempty"`
	AdCb                  string          `json:"ad_cb,omitempty"`
	ResourceID            int             `json:"resource_id,omitempty"`
	RequestID             string          `json:"request_id,omitempty"`
	ClientIP              string          `json:"client_ip,omitempty"`
	IsAd                  bool            `json:"is_ad"`
	IsAdLoc               bool            `json:"is_ad_loc,omitempty"`
	Schema                string          `json:"schema,omitempty"`
	SchemaTitle           string          `json:"schema_title,omitempty"`
	SchemaPackageName     string          `json:"schema_package_name,omitempty"`
	SchemaCallupWhiteList []string        `json:"schema_callup_white_list,omitempty"`
	Extra                 json.RawMessage `json:"extra,omitempty"`
}

type Show struct {
	ID    int64      `json:"id"`
	Stime xtime.Time `json:"stime"`
	Etime xtime.Time `json:"etime"`
}

type CmSplash struct {
	*CmConfig
	List []*List `json:"list,omitempty"`
	Show []*Show `json:"show,omitempty"`
}

type CmConfig struct {
	MaxTime      int `json:"max_time"`
	MinInterval  int `json:"min_interval"`
	PullInterval int `json:"pull_interval"`
}

// PlatChange
func (s *Splash) PlatChange() {
	switch s.Plat {
	case 1: // resource iphone
		s.Plat = model.PlatIPhone
	case 2: // resource android
		s.Plat = model.PlatAndroid
	case 3: // resource pad
		s.Plat = model.PlatIPad
	case 4: // resource iphoneg
		s.Plat = model.PlatIPhoneI
	case 5: // resource androidg
		s.Plat = model.PlatAndroidG
	case 6: // resource padg
		s.Plat = model.PlatIPadI
	case 8: // resource androidi
		s.Plat = model.PlatAndroidI
	}
	if s.Operate == 1 { // NOTE: operate=1 means AD
		s.Type = 1 // NOTE: type=1 compatiable mobile, must type=1 can splash
	}
}

// BirthDate
func (s *Splash) BirthDate() {
	s.BirthStart = s.Start.Time().Format("0102")
	s.BirthEnd = s.End.Time().Format("0102")
	s.BirthStartMonth = s.Start.Time().Format("01")
	s.BirthEndMonth = s.End.Time().Format("01")
	s.Start = xtime.Time(0)
	s.End = xtime.Time(0)
}

// Ratio calc width/height ratio.
func Ratio(w, h int) float64 {
	return float64(w) / float64(h)
}
