package operate

import (
	"encoding/json"
	"sort"
	"strconv"

	"go-common/app/interface/main/app-card/model"
	"go-common/library/log"
	xtime "go-common/library/time"
)

type Menu struct {
	TabID       int64               `json:"tab_id,omitempty"`
	Name        string              `json:"name,omitempty"`
	Img         string              `json:"img,omitempty"`
	Icon        string              `json:"icon,omitempty"`
	Color       string              `json:"color,omitempty"`
	ID          int64               `json:"id,omitempty"`
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
	PlatStr   string `json:"plat,omitempty"`
	BuildStr  string `json:"build,omitempty"`
	Condition string `json:"conditions,omitempty"`
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

type Active struct {
	ID         int64           `json:"id,omitempty"`
	ParentID   int64           `json:"parent_id,omitempty"`
	Name       string          `json:"name,omitempty"`
	Background string          `json:"background,omitempty"`
	Type       string          `json:"type,omitempty"`
	Content    json.RawMessage `json:"content,omitempty"`
	// Extra
	Pid      int64     `json:"pid,omitempty"`
	Title    string    `json:"title,omitempty"`
	Subtitle string    `json:"subtitle,omitempty"`
	Desc     string    `json:"desc,omitempty"`
	Param    string    `json:"param,omitempty"`
	Goto     model.Gt  `json:"goto,omitempty"`
	Cover    string    `json:"cover,omitempty"`
	Limit    int       `json:"limit,omitempty"`
	Items    []*Active `json:"items,omitempty"`
}

type Actives []*Active

func (is Actives) Len() int           { return len(is) }
func (is Actives) Less(i, j int) bool { return is[i].ID < is[j].ID }
func (is Actives) Swap(i, j int)      { is[i], is[j] = is[j], is[i] }

type CardItem struct {
	Title     string `json:"title,omitempty"`
	Cover     string `json:"cover,omitempty"`
	LinkType  string `json:"link_type,omitempty"`
	LinkValue string `json:"link_value,omitempty"`
	Weight    int    `json:"weight,omitempty"`
}

type CardItems []*CardItem

func (is CardItems) Len() int           { return len(is) }
func (is CardItems) Less(i, j int) bool { return is[i].Weight < is[j].Weight }
func (is CardItems) Swap(i, j int)      { is[i], is[j] = is[j], is[i] }

func (a *Active) Change() {
	switch a.Type {
	case "archive":
		a.Type = "player"
	case "live":
		a.Type = "player_live"
	case "basic":
		a.Type = "content_rcmd"
	case "shortcut":
		a.Type = "entrance"
	case "common":
		a.Type = "background"
	case "tag":
		a.Type = "tag_rcmd"
	}
	switch a.Type {
	// 基本类型
	case "player_live", "converge", "special", "article_s", "player":
		var id int64
		if err := json.Unmarshal(a.Content, &id); err != nil {
			log.Error("%+v", err)
			return
		}
		if id > 0 {
			a.Pid = id
		}
	// 新增类型
	case "content_rcmd":
		var basic struct {
			Type     string `json:"type,omitempty"`
			Title    string `json:"title,omitempty"`
			Subtitle string `json:"subtitle,omitempty"`
			Sublink  string `json:"sublink,omitempty"`
			Content  []*struct {
				LinkType  string `json:"link_type,omitempty"`
				LinkValue string `json:"link_value,omitempty"`
			} `json:"content,omitempty"`
		}
		if err := json.Unmarshal(a.Content, &basic); err != nil {
			log.Error("%+v", err)
			return
		}
		ris := make([]*Active, 0, len(basic.Content))
		for _, c := range basic.Content {
			typ, _ := strconv.Atoi(c.LinkType)
			id, _ := strconv.ParseInt(c.LinkValue, 10, 64)
			ri := &Active{Goto: model.OperateType[typ], Pid: id}
			if ri.Goto != "" {
				ris = append(ris, ri)
			}
		}
		if len(ris) != 0 {
			a.Items = ris
			a.Title = basic.Title
			a.Subtitle = basic.Subtitle
			a.Param = basic.Sublink
		}
	case "entrance", "banner":
		var card struct {
			Type     string      `json:"type,omitempty"`
			CardItem []*CardItem `json:"card_item,omitempty"`
		}
		if err := json.Unmarshal(a.Content, &card); err != nil {
			log.Error("%+v", err)
			return
		}
		ris := make([]*Active, 0, len(card.CardItem))
		sort.Sort(CardItems(card.CardItem))
		for _, item := range card.CardItem {
			typ, _ := strconv.Atoi(item.LinkType)
			id, _ := strconv.ParseInt(item.LinkValue, 10, 64)
			ri := &Active{Goto: model.OperateType[typ], Pid: id, Param: item.LinkValue, Title: item.Title, Cover: item.Cover}
			if ri.Goto != "" {
				ris = append(ris, ri)
			}
		}
		if len(ris) != 0 {
			a.Items = ris
		}
	case "tag_rcmd":
		var tag struct {
			AidStr    string `json:"aid,omitempty"`
			Type      string `json:"type,omitempty"`
			NumberStr string `json:"number,omitempty"`
			Tid       int64  `json:"-"`
			Number    int    `json:"-"`
		}
		if err := json.Unmarshal(a.Content, &tag); err != nil {
			log.Error("%+v", err)
			return
		}
		tag.Tid, _ = strconv.ParseInt(tag.AidStr, 10, 64)
		tag.Number, _ = strconv.Atoi(tag.NumberStr)
		if tag.Tid != 0 {
			a.Pid = tag.Tid
			a.Limit = tag.Number
		}
	case "background":
		a.Title = a.Name
		a.Cover = a.Background
	case "news":
		var news struct {
			Title string `json:"title,omitempty"`
			Body  string `json:"body,omitempty"`
			Link  string `json:"link,omitempty"`
		}
		if err := json.Unmarshal(a.Content, &news); err != nil {
			log.Error("%+v", err)
			return
		}
		if news.Body != "" {
			a.Title = news.Title
			a.Desc = news.Body
			a.Param = news.Link
		}
	}
}
