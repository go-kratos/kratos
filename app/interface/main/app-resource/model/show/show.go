package show

import (
	"go-common/app/interface/main/app-resource/model"
	"go-common/app/interface/main/app-resource/model/tab"
	resource "go-common/app/service/main/resource/model"
	"strconv"
)

type Tab struct {
	ID              int64  `json:"id,omitempty"`
	Icon            string `json:"icon,omitempty"`
	IconSelected    string `json:"icon_selected,omitempty"`
	Name            string `json:"name,omitempty"`
	URI             string `json:"uri,omitempty"`
	TabID           string `json:"tab_id,omitempty"`
	Color           string `json:"color,omitempty"`
	Pos             int    `json:"pos,omitempty"`
	DefaultSelected int    `json:"default_selected,omitempty"`
	Module          int    `json:"-"`
	ModuleStr       string `json:"-"`
	Plat            int8   `json:"-"`
	Group           string `json:"-"`
	Language        string `json:"-"`
}

type Limit struct {
	ID        int64  `json:"-"`
	Build     int    `json:"-"`
	Condition string `json:"-"`
}

func (t *Tab) TabChange(rsb *resource.SideBar, abtest map[string]string, defaultTab map[string]*Tab) (ok bool) {
	var (
		_top    = 10
		_tab    = 8
		_bottom = 9
	)
	t.ID = rsb.ID
	t.Icon = rsb.Logo
	t.IconSelected = rsb.LogoSelected
	t.Name = rsb.Name
	t.URI = rsb.Param
	t.Module = rsb.Module
	t.Plat = rsb.Plat
	t.Language = rsb.Language
	switch t.Module {
	case _top:
		t.ModuleStr = "top"
	case _tab:
		t.ModuleStr = "tab"
		t.Icon = ""
		t.IconSelected = ""
	case _bottom:
		t.ModuleStr = "bottom"
	default:
		return false
	}
	if len(abtest) > 0 {
		if groups, ok := abtest[t.URI]; ok {
			t.Group = groups
		}
	}
	if len(defaultTab) > 0 {
		if dt, ok := defaultTab[t.URI]; ok && dt != nil {
			t.DefaultSelected = dt.DefaultSelected
			t.TabID = dt.TabID
		}
		if rsb.TabID != "" {
			t.TabID = rsb.TabID
		}
	}
	return true
}

func (t *Tab) TabMenuChange(m *tab.Menu) {
	t.TabID = strconv.FormatInt(m.TabID, 10)
	t.Name = m.Name
	t.Color = m.Color
	t.ID = m.ID
	t.ModuleStr = "tab"
	t.URI = model.FillURI(model.GotoPegasusTab, strconv.FormatInt(t.ID, 10), model.PegasusHandler(m))
}
