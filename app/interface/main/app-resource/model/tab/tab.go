package tab

import (
	"encoding/json"
	"strconv"

	"go-common/library/log"
	xtime "go-common/library/time"
)

type Menu struct {
	ID          int64               `json:"id,omitempty"`
	Name        string              `json:"name,omitempty"`
	Img         string              `json:"img,omitempty"`
	Icon        string              `json:"icon,omitempty"`
	Color       string              `json:"color,omitempty"`
	TabID       int64               `json:"tab_id,omitempty"`
	Plat        int                 `json:"-"`
	CType       int                 `json:"-"`
	CValue      string              `json:"-"`
	PlatVersion json.RawMessage     `json:"-"`
	STime       xtime.Time          `json:"-"`
	ETime       xtime.Time          `json:"-"`
	Status      int                 `json:"-"`
	Badge       string              `json:"-"`
	Versions    map[int8][]*Version `json:"-"`
}

type Version struct {
	PlatStr   string `json:"plat"`
	BuildStr  string `json:"build"`
	Condition string `json:"conditions"`
	Plat      int8   `json:"-"`
	Build     int    `json:"-"`
}

func (m *Menu) Change() {
	m.Icon = m.Badge
	var vs []*Version
	if err := json.Unmarshal(m.PlatVersion, &vs); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", m.PlatVersion, err)
		return
	}
	vm := make(map[int8][]*Version, len(vs))
	for _, v := range vs {
		if v.PlatStr == "" || v.BuildStr == "" {
			continue
		}
		if plat, err := strconv.ParseInt(v.PlatStr, 10, 8); err != nil {
			log.Error("strconv.ParseInt(%s,10,8) error(%v)", v.PlatStr, err)
			continue
		} else {
			v.Plat = int8(plat)
		}
		if build, err := strconv.Atoi(v.BuildStr); err != nil {
			log.Error("strconv.Atoi(%s) error(%v)", v.BuildStr, err)
			continue
		} else {
			v.Build = build
		}
		vm[v.Plat] = append(vm[v.Plat], v)
	}
	m.Versions = vm
	if m.CType == 1 {
		var err error
		if m.ID, err = strconv.ParseInt(m.CValue, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", m.CValue, err)
			return
		}
	}
}
