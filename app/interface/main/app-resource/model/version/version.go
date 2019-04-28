package version

import (
	xtime "go-common/library/time"
	"strings"
)

const (
	PlatAndroid  = int8(0)
	PlatIPhone   = int8(1)
	PlatIPad     = int8(2)
	PlatWinPhone = int8(3)
)

type Version struct {
	Id      int        `json:"-"`
	Plat    int8       `json:"plat"`
	Desc    string     `json:"desc"`
	Version string     `json:"version"`
	Build   int        `json:"build"`
	PTime   xtime.Time `json:"ptime"`
}

type VersionUpdate struct {
	Id         int                 `json:"-"`
	Channel    string              `json:"-"`
	Coverage   int                 `json:"-"`
	Version    string              `json:"ver"`
	Build      int                 `json:"build"`
	Desc       string              `json:"info"`
	State      int                 `json:"-"`
	Size       string              `json:"size"`
	Url        string              `json:"url"`
	MD5        string              `json:"hash"`
	SdkInts    string              `json:"-"`
	SdkIntList map[string]struct{} `json:"-"`
	Model      string              `json:"-"`
	Policy     int                 `json:"policy"`
	Plat       int8                `json:"-"`
	IsForce    int                 `json:"is_force"`
	IsPush     int                 `json:"is_push"`
	PolicyName string              `json:"-"`
	IsGray     int                 `json:"is_gray"`
	PolicyURL  string              `json:"policy_url,omitempty"`
	BuvidStart int                 `json:"-"`
	BuvidEnd   int                 `json:"-"`
	Mtime      xtime.Time          `json:"mtime"`
	Incre      *Incremental        `json:"patch,omitempty"`
}

type UpdateLimit struct {
	ID         int    `json:"-"`
	BuildLimit int    `json:"-"`
	Conditions string `json:"-"`
}

type VersionSo struct {
	Id           int    `json:"-"`
	Package      string `json:"-"`
	Name         string `json:"-"`
	Description  string `json:"-"`
	Clear        int    `json:"-"`
	Ver_code     int    `json:"ver_code"`
	Ver_name     string `json:"ver_name"`
	Url          string `json:"url"`
	Size         int    `json:"size"`
	Enable_state int    `json:"enable"`
	Force_state  int    `json:"force"`
	Md5          string `json:"md5"`
	Min_build    int    `json:"min_build"`
	Coverage     int    `json:"-"`
	Sdkint       int    `json:"-"`
	Model        string `json:"-"`
}

type VersionSoDesc struct {
	Package     string       `json:"package"`
	Name        string       `json:"name"`
	Description string       `json:"desc"`
	Clear       int          `json:"clear"`
	Versions    []*VersionSo `json:"versions"`
}

// Incremental version Incremental
type Incremental struct {
	ID            int    `json:"-"`
	TargetVersion string `json:"-"`
	TargetBuild   int    `json:"-"`
	TargetID      string `json:"new_id"`
	SourceVersion string `json:"-"`
	SourceBuild   int    `json:"-"`
	SourceID      string `json:"old_id"`
	TaskID        string `json:"-"`
	FilePath      string `json:"-"`
	URL           string `json:"url"`
	Md5           string `json:"md5"`
	Size          int    `json:"size"`
	Policy        int    `json:"-"`
	Plat          int8   `json:"-"`
	Build         int    `json:"-"`
}

// Rn
type Rn struct {
	ID            int    `json:"-"`
	DeploymentKey string `json:"-"`
	BundleID      string `json:"bundle_id"`
	URL           string `json:"url"`
	Md5           string `json:"md5"`
	Size          int    `json:"size"`
	Version       string `json:"-"`
}

// VersionUpdateChange version update change
func (v *VersionUpdate) VersionUpdateChange() {
	if v.SdkInts != "" {
		v.SdkIntList = map[string]struct{}{}
		tmp := strings.Split(v.SdkInts, ",")
		for _, sdkint := range tmp {
			v.SdkIntList[sdkint] = struct{}{}
		}
	}
}
