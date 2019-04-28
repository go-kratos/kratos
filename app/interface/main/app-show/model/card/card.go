package card

import (
	"encoding/json"
	"strconv"

	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/app/interface/main/app-show/model"
	xtime "go-common/library/time"
)

type Column struct {
	ID        int    `json:"-"`
	Tab       int    `json:"-"`
	RegionID  int    `json:"-"`
	Tpl       int    `json:"-"`
	Name      string `json:"-"`
	Desc      string `json:"-"`
	PlatVer   string `json:"-"`
	Plat      int8   `json:"-"`
	Build     int    `json:"-"`
	Condition string `json:"-"`
	Type      string `json:"-"`
}

type ColumnNper struct {
	ID        int        `json:"-"`
	ColumnID  int        `json:"-"`
	Name      string     `json:"-"`
	Desc      string     `json:"-"`
	Nper      string     `json:"-"`
	NperTime  xtime.Time `json:"-"`
	Cover     string     `json:"-"`
	PlatVer   string     `json:"-"`
	Title     string     `json:"-"`
	Rtype     int        `json:"-"`
	Rvalue    string     `json:"-"`
	Plat      int8       `json:"-"`
	Build     int        `json:"-"`
	Condition string     `json:"-"`
	Goto      string     `json:"-"`
	Param     string     `json:"-"`
	URI       string     `json:"-"`
}

type Content struct {
	ID     int    `json:"-"`
	Module int    `json:"-"`
	RecID  int    `json:"-"`
	Type   int8   `json:"-"`
	Value  string `json:"-"`
	Title  string `json:"-"`
	TagID  int    `json:"-"`
}

type Card struct {
	ID        int    `json:"-"`
	Tab       int    `json:"-"`
	RegionID  int    `json:"-"`
	Type      int    `json:"-"`
	Title     string `json:"-"`
	Cover     string `json:"-"`
	Rtype     int    `json:"-"`
	Rvalue    string `json:"-"`
	PlatVer   string `json:"-"`
	Plat      int8   `json:"-"`
	Build     int    `json:"-"`
	Condition string `json:"-"`
	TypeStr   string `json:"-"`
	Goto      string `json:"-"`
	Param     string `json:"-"`
	URi       string `json:"-"`
	Desc      string `json:"-"`
	TagID     int    `json:"-"`
}

type ColumnList struct {
	Cid       int    `json:"-"`
	Ceid      int    `json:"-"`
	Name      string `json:"-"`
	Cname     string `json:"-"`
	PlatVer   string `json:"-"`
	Plat      int8   `json:"-"`
	Build     int    `json:"-"`
	Condition string `json:"-"`
}

type PlatLimit struct {
	Plat      int8   `json:"plat"`
	Build     int    `json:"build"`
	Condition string `json:"conditions"`
}

type PopularCard struct {
	ID              int64                       `json:"id"`
	Title           string                      `json:"title"`
	ChannelID       int64                       `json:"channel_id"`
	Type            string                      `json:"type"`
	Value           int64                       `json:"value"`
	Reason          string                      `json:"reason"`
	ReasonType      int8                        `json:"reason_type"`
	Pos             int                         `json:"pos"`
	FromType        string                      `json:"from_type"`
	PopularCardPlat map[int8][]*PopularCardPlat `json:"popularcardplat"`
	Idx             int64                       `json:"-"`
	CornerMark      int8                        `json:"corner_mark"`
}

type PopularCardPlat struct {
	CardID    int64  `json:"card_id"`
	Plat      int8   `json:"plat"`
	Condition string `json:"condition"`
	Build     int    `json:"build"`
}

func (c *Card) CardPlatChange() (platlinits []*PlatLimit) {
	platlinits = platJsonChange(c.PlatVer)
	return
}

func (c *Column) ColumnPlatChange() (platlinits []*PlatLimit) {
	platlinits = platJsonChange(c.PlatVer)
	return
}

func (c *ColumnList) ColumnListPlatChange() (platlinits []*PlatLimit) {
	platlinits = platJsonChange(c.PlatVer)
	return
}

func (c *ColumnNper) ColumnNperPlatChange() (platlinits []*PlatLimit) {
	platlinits = platJsonChange(c.PlatVer)
	return
}

func (c *Card) CardGotoChannge() {
	c.TypeStr = cardTypeChange(c.Type)
	c.Goto, c.Param, c.URi = gotoURI(c.Rtype, c.Rvalue)
}

func (c *Column) ColumnGotoChannge() {
	switch c.Tpl {
	case 1:
		c.Type = model.GotoDaily
	case 2:
		c.Type = model.GotoColumn
	}
}

func (c *ColumnNper) ColumnNperGotoChange() {
	c.Goto, c.Param, c.URI = gotoURI(c.Rtype, c.Rvalue)
}

func gotoURI(typeInt int, value string) (gotoStr, paramStr, uri string) {
	switch typeInt {
	case 1:
		gotoStr = model.GotoDaily
	case 4:
		gotoStr = model.GotoWeb
	case 5:
		gotoStr = model.GotoAv
	case 6:
		gotoStr = model.GotoLive
	case 7:
		gotoStr = model.GotoBangumi
	case 8:
		gotoStr = model.GotoGame
	case 9:
		gotoStr = model.GotoColumn
	case 10:
		gotoStr = model.GotoColumnStage
	case 11:
		gotoStr = model.GotoArticle
	default:
		return
	}
	paramStr = value
	uri = model.FillURI(gotoStr, paramStr, nil)
	return
}

func cardTypeChange(cardInt int) (cardStr string) {
	switch cardInt {
	case 1:
		cardStr = model.GotoDaily
	case 2:
		cardStr = model.GotoTopic
	case 3:
		cardStr = model.GotoActivity
	case 4:
		cardStr = model.GotoRank
	case 5:
		cardStr = model.GotoCard
	case 6:
		cardStr = model.GotoVeidoCard
	case 7:
		cardStr = model.GotoSpecialCard
	case 8:
		cardStr = model.GotoTagCard
	}
	return
}

// platJsonChange json change plat build condition
func platJsonChange(jsonStr string) (platlinits []*PlatLimit) {
	var tmp []struct {
		Plat      string `json:"plat"`
		Build     string `json:"build"`
		Condition string `json:"conditions"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &tmp); err == nil {
		for _, limit := range tmp {
			platlinit := &PlatLimit{}
			switch limit.Plat {
			case "0": // resource android
				platlinit.Plat = model.PlatAndroid
			case "1": // resource iphone
				platlinit.Plat = model.PlatIPhone
			case "2": // resource pad
				platlinit.Plat = model.PlatIPad
			case "5": // resource iphone_i
				platlinit.Plat = model.PlatIPhoneI
			case "8": // resource android_i
				platlinit.Plat = model.PlatAndroidI
			}
			platlinit.Build, _ = strconv.Atoi(limit.Build)
			platlinit.Condition = limit.Condition
			platlinits = append(platlinits, platlinit)
		}
	}
	return
}

func (c *PopularCard) PopularCardToAiChange() (a *ai.Item) {
	a = &ai.Item{
		Goto:       c.Type,
		ID:         c.Value,
		RcmdReason: c.fromRcmdReason(),
		CornerMark: c.CornerMark,
	}
	if c.CornerMark == 2 {
		a.CornerMark = 4
	}
	return
}

func (c *PopularCard) fromRcmdReason() (a *ai.RcmdReason) {
	var content string
	switch c.ReasonType {
	case 0:
		content = ""
	case 1:
		content = "编辑精选"
	case 2:
		content = "热门推荐"
	case 3:
		content = c.Reason
	}
	if content != "" {
		a = &ai.RcmdReason{ID: 1, Content: content, BgColor: "yellow", IconLocation: "left_top"}
	}
	return
}
