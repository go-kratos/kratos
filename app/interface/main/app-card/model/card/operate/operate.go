package operate

import (
	"encoding/json"
	"go-common/app/interface/main/app-card/model"
	"go-common/app/interface/main/app-card/model/card/banner"
	"go-common/app/interface/main/app-card/model/card/cm"
	"go-common/app/interface/main/app-card/model/card/rank"
	"go-common/library/log"
	"sort"
	"strconv"
)

type Card struct {
	Plat       int8                          `json:"plat,omitempty"`
	Build      int                           `json:"build,omitempty"`
	ID         int64                         `json:"id,omitempty"`
	Param      string                        `json:"param,omitempty"`
	CardGoto   model.CardGt                  `json:"card_goto,omitempty"`
	Goto       model.Gt                      `json:"goto,omitempty"`
	URI        string                        `json:"uri,omitempty"`
	Title      string                        `json:"title,omitempty"`
	Desc       string                        `json:"desc,omitempty"`
	Cover      string                        `json:"cover,omitempty"`
	Coverm     map[model.ColumnStatus]string `json:"coverm,omitempty"`
	Avatar     string                        `json:"avatar,omitempty"`
	Download   int32                         `json:"download,omitempty"`
	Badge      string                        `json:"badge,omitempty"`
	Ratio      int                           `json:"ratio,omitempty"`
	Score      int32                         `json:"score,omitempty"`
	Tid        int64                         `json:"tid,omitempty"`
	Subtitle   string                        `json:"subtitle,omitempty"`
	Limit      int                           `json:"limit,omitempty"`
	Items      []*Card                       `json:"items,omitempty"`
	AdInfo     *cm.AdInfo                    `json:"ad_info,omitempty"`
	Banner     []*banner.Banner              `json:"banner,omitempty"`
	Hash       string                        `json:"verson,omitempty"`
	TrackID    string                        `json:"trackid,omitempty"`
	FromType   string                        `json:"from_type,omitempty"`
	ShowUGCPay bool                          `json:"show_ucg_pay,omitempty"`
	Switch     model.Switch                  `json:"switch,omitempty"`
	SwitchLike model.Switch                  `json:"switch_like,omitempty"`
	Buttons    []*Button                     `json:"buttons,omitempty"`
}

type Button struct {
	Text  string `json:"text,omitempty"`
	Event string `json:"event,omitempty"`
}

func (c *Card) From(cardGoto model.CardGt, id int64, tid int64, plat int8, build int) {
	c.CardGoto = cardGoto
	c.ID = id
	c.Tid = tid
	c.Goto = model.Gt(cardGoto)
	c.Param = strconv.FormatInt(id, 10)
	c.URI = strconv.FormatInt(id, 10)
	c.Plat = plat
	c.Build = build
}

func (c *Card) FromSwitch(sw model.Switch) {
	c.SwitchLike = sw
}

func (c *Card) FromDownload(o *Download) {
	c.CardGoto = model.CardGotoDownload
	c.Param = strconv.FormatInt(o.ID, 10)
	c.Coverm = map[model.ColumnStatus]string{model.ColumnSvrSingle: o.Cover, model.ColumnSvrDouble: o.DoubleCover}
	c.Title = o.Title
	c.Goto = model.OperateType[o.URLType]
	c.URI = o.URLValue
	c.Avatar = o.Icon
	c.Download = o.Number
	c.Desc = o.Desc
}

func (c *Card) FromSpecial(o *Special) {
	c.CardGoto = model.CardGotoSpecial
	c.Param = strconv.FormatInt(o.ID, 10)
	c.Coverm = map[model.ColumnStatus]string{model.ColumnSvrSingle: o.SingleCover, model.ColumnSvrDouble: o.Cover}
	c.Title = o.Title
	c.Goto = model.OperateType[o.ReType]
	c.URI = o.ReValue
	c.Desc = o.Desc
	c.Badge = o.Badge
	if o.Size == "1020x300" {
		c.Ratio = 34
	} else if o.Size == "1020x378" {
		c.Ratio = 27
	}
}

func (c *Card) FromTopstick(o *Special) {
	c.CardGoto = model.CardGotoTopstick
	c.Param = strconv.FormatInt(o.ID, 10)
	c.Title = o.Title
	c.Goto = model.OperateType[o.ReType]
	c.URI = o.ReValue
	c.Desc = o.Desc
	c.Badge = o.Badge
}

func (c *Card) FromFollow(o *Follow) {
	switch o.Type {
	case "upper", "channel_three":
		var contents []*struct {
			Ctype  string `json:"ctype,omitempty"`
			Cvalue int64  `json:"cvalue,omitempty"`
		}
		if err := json.Unmarshal(o.Content, &contents); err != nil {
			log.Error("%+v", err)
			return
		}
		items := make([]*Card, 0, len(contents))
		for _, content := range contents {
			var gt model.Gt
			switch content.Ctype {
			case "mid":
				gt = model.GotoMid
			case "channel_id":
				gt = model.GotoTag
			default:
				continue
			}
			items = append(items, &Card{ID: content.Cvalue, Goto: gt, Param: strconv.FormatInt(content.Cvalue, 10), URI: strconv.FormatInt(content.Cvalue, 10)})
		}
		if len(items) < 3 {
			return
		}
		c.Items = items
		c.CardGoto = model.CardGotoSubscribe
		c.Title = o.Title
		c.Param = strconv.FormatInt(o.ID, 10)
	case "channel_single":
		var content struct {
			Aid       int64 `json:"aid"`
			ChannelID int64 `json:"channel_id"`
		}
		if err := json.Unmarshal(o.Content, &content); err != nil {
			log.Error("%+v", err)
			return
		}
		c.CardGoto = model.CardGotoChannelRcmd
		c.Title = o.Title
		c.ID = content.Aid
		c.Tid = content.ChannelID
		c.Goto = model.GotoAv
		c.Param = strconv.FormatInt(o.ID, 10)
		c.URI = strconv.FormatInt(content.Aid, 10)
	}
}

func (c *Card) FromConverge(o *Converge) {
	c.CardGoto = model.CardGotoConverge
	c.Param = strconv.FormatInt(o.ID, 10)
	c.Coverm = map[model.ColumnStatus]string{model.ColumnSvrSingle: o.Cover, model.ColumnSvrDouble: o.Cover}
	c.Title = o.Title
	c.Goto = model.OperateType[o.ReType]
	c.URI = o.ReValue
	var contents []*struct {
		Ctype  string `json:"ctype,omitempty"`
		Cvalue string `json:"cvalue,omitempty"`
	}
	if err := json.Unmarshal(o.Content, &contents); err != nil {
		log.Error("%+v", err)
		return
	}
	c.Items = make([]*Card, 0, len(contents))
	for _, content := range contents {
		var (
			gt     model.Gt
			cardGt model.CardGt
		)
		id, _ := strconv.ParseInt(content.Cvalue, 10, 64)
		if id == 0 {
			continue
		}
		switch content.Ctype {
		case "0":
			gt = model.GotoAv
			cardGt = model.CardGotoAv
		case "1":
			gt = model.GotoLive
			cardGt = model.CardGotoLive
		case "2":
			gt = model.GotoArticle
			cardGt = model.CardGotoArticleS
		default:
			continue
		}
		c.Items = append(c.Items, &Card{ID: id, CardGoto: cardGt, Goto: gt, Param: content.Cvalue, URI: content.Cvalue})
	}
}

func (c *Card) FromRank(os []*rank.Rank) {
	c.CardGoto = model.CardGotoRank
	c.Goto = model.GotoRank
	c.Items = make([]*Card, 0, len(os))
	for _, o := range os {
		c.Items = append(c.Items, &Card{Goto: model.GotoAv, ID: o.Aid, Param: strconv.FormatInt(o.Aid, 10), URI: strconv.FormatInt(o.Aid, 10), Score: o.Score})
	}
}

func (c *Card) FromActive(o *Active) {
	switch o.Type {
	case "live", "player_live", "converge", "special", "archive", "player":
		var id int64
		if err := json.Unmarshal(o.Content, &id); err != nil {
			log.Error("%+v", err)
			return
		}
		if id < 1 {
			return
		}
		c.ID = id
		c.Param = strconv.FormatInt(id, 10)
		switch o.Type {
		case "live":
			c.CardGoto = model.CardGotoPlayerLive
		case "converge":
			c.CardGoto = model.CardGotoConverge
		case "special":
			c.CardGoto = model.CardGotoSpecial
		case "archive":
			c.CardGoto = model.CardGotoPlayer
		}
	case "basic", "content_rcmd":
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
		if err := json.Unmarshal(o.Content, &basic); err != nil {
			log.Error("%+v", err)
			return
		}
		items := make([]*Card, 0, len(basic.Content))
		for _, c := range basic.Content {
			typ, _ := strconv.Atoi(c.LinkType)
			id, _ := strconv.ParseInt(c.LinkValue, 10, 64)
			ri := &Card{Goto: model.OperateType[typ], ID: id, Param: c.LinkValue}
			if ri.Goto != "" {
				items = append(items, ri)
			}
		}
		if len(items) == 0 {
			return
		}
		c.Items = items
		c.Title = basic.Title
		c.Subtitle = basic.Subtitle
		c.URI = basic.Sublink
		c.CardGoto = model.CardGotoContentRcmd
	case "shortcut", "entrance", "banner":
		var card struct {
			Type     string      `json:"type,omitempty"`
			CardItem []*CardItem `json:"card_item,omitempty"`
		}
		if err := json.Unmarshal(o.Content, &card); err != nil {
			log.Error("%+v", err)
			return
		}
		items := make([]*Card, 0, len(card.CardItem))
		sort.Sort(CardItems(card.CardItem))
		for _, v := range card.CardItem {
			typ, _ := strconv.Atoi(v.LinkType)
			id, _ := strconv.ParseInt(v.LinkValue, 10, 64)
			item := &Card{Goto: model.OperateType[typ], ID: id, Param: v.LinkValue, URI: v.LinkValue, Title: v.Title, Cover: v.Cover}
			if item.Goto != "" {
				items = append(items, item)
			}
		}
		if len(items) == 0 {
			return
		}
		c.Items = items
		switch o.Type {
		case "shortcut", "entrance":
			c.CardGoto = model.CardGotoEntrance
		case "banner":
			c.CardGoto = model.CardGotoBanner
		}
	case "common", "background":
		c.Title = o.Name
		c.Cover = o.Background
	case "tag", "tag_rcmd":
		var tag struct {
			AidStr    string `json:"aid,omitempty"`
			Type      string `json:"type,omitempty"`
			NumberStr string `json:"number,omitempty"`
			Tid       int64  `json:"-"`
			Number    int    `json:"-"`
		}
		if err := json.Unmarshal(o.Content, &tag); err != nil {
			log.Error("%+v", err)
			return
		}
		tag.Tid, _ = strconv.ParseInt(tag.AidStr, 10, 64)
		tag.Number, _ = strconv.Atoi(tag.NumberStr)
		if tag.Tid == 0 {
			return
		}
		c.ID = tag.Tid
		c.Limit = tag.Number
		c.Goto = model.GotoTag
		c.CardGoto = model.CardGotoTagRcmd
		c.Subtitle = "查看更多"
	case "news":
		var news struct {
			Title string `json:"title,omitempty"`
			Body  string `json:"body,omitempty"`
			Link  string `json:"link,omitempty"`
		}
		if err := json.Unmarshal(o.Content, &news); err != nil {
			log.Error("%+v", err)
			return
		}
		if news.Body == "" {
			return
		}
		c.Title = news.Title
		c.Desc = news.Body
		c.URI = news.Link
		c.Goto = model.GotoWeb
		c.CardGoto = model.CardGotoNews
	}
	c.Title = o.Title
	c.Param = strconv.FormatInt(o.ID, 10)
}

func (c *Card) FromAdAv(o *cm.AdInfo) {
	c.CardGoto = model.CardGotoAdAv
	c.AdInfo = o
}

func (c *Card) FromActiveBanner(os []*Active, hash string) {
	c.Banner = make([]*banner.Banner, 0, len(os))
	for _, o := range os {
		banner := &banner.Banner{ID: o.Pid, Title: o.Title, Image: o.Cover, URI: model.FillURI(o.Goto, o.Param, nil)}
		c.Banner = append(c.Banner, banner)
	}
	c.CardGoto = model.CardGotoBanner
	c.Hash = hash
}

func (c *Card) FromBanner(os []*banner.Banner, hash string) {
	if len(os) == 0 {
		return
	}
	c.Banner = os
	c.CardGoto = model.CardGotoBanner
	c.Hash = hash
}

func (c *Card) FromLogin(o int64) {
	if !model.IsIPad(c.Plat) {
		if o != 0 {
			c.Param = strconv.FormatInt(o, 10)
		} else {
			c.Param = "1"
		}
	} else {
		c.Param = "5"
	}
	c.CardGoto = model.CardGotoLogin
}

func (c *Card) FromCardSet(o *CardSet) {
	switch o.Type {
	case "pgcs_rcmd":
		var contents []*struct {
			ID interface{} `json:"id,omitempty"`
		}
		if err := json.Unmarshal(o.Content, &contents); err != nil {
			log.Error("%+v", err)
			return
		}
		for _, content := range contents {
			var cid int64
			switch v := content.ID.(type) {
			case string:
				cid, _ = strconv.ParseInt(v, 10, 64)
			case float64:
				cid = int64(v)
			}
			item := &Card{ID: cid, Goto: model.GotoPGC}
			c.Items = append(c.Items, item)
		}
		c.Title = o.Title
		c.Param = strconv.FormatInt(o.ID, 10)
		c.CardGoto = model.CardGotoPgcsRcmd
	case "up_rcmd_new":
		var contents []*struct {
			ID interface{} `json:"id,omitempty"`
		}
		if err := json.Unmarshal(o.Content, &contents); err != nil {
			log.Error("%+v", err)
			return
		}
		for _, content := range contents {
			var aid int64
			switch v := content.ID.(type) {
			case string:
				aid, _ = strconv.ParseInt(v, 10, 64)
			case float64:
				aid = int64(v)
			}
			item := &Card{ID: aid, Goto: model.GotoAv}
			c.Items = append(c.Items, item)
		}
		c.Title = "新星卡片"
		c.Desc = o.Title
		c.Param = strconv.FormatInt(o.Value, 10)
		c.ID = o.Value
		c.CardGoto = model.CardGotoUpRcmdNew
	}
}

func (c *Card) FromFollowMode(title, desc string, button []string) {
	c.Title = title
	if c.Title == "" {
		c.Title = "启用首页推荐 - 关注模式（内测版）"
	}
	c.Desc = desc
	if c.Desc == "" {
		c.Desc = "我们根据你对bilibili推荐的反馈，为你定制了关注模式。开启后，仅为你显示关注UP主更新的视频哦。尝试体验一下？"
	}
	if len(button) == 2 {
		c.Buttons = []*Button{
			{Text: button[0], Event: "close"},
			{Text: button[1], Event: "follow_mode"},
		}
	} else {
		c.Buttons = []*Button{
			{Text: "暂不需要", Event: "close"},
			{Text: "立即开启", Event: "follow_mode"},
		}
	}
	c.CardGoto = model.CardGotoFollowMode
}

func (c *Card) FromEventTopic(o *EventTopic) {
	c.Title = o.Title
	c.Desc = o.Desc
	c.Cover = o.Cover
	switch o.ReType {
	case 1:
		c.Goto = model.Gt("topic")
	case 2:
		c.Goto = model.Gt("broadcast")
	case 3:
		c.Goto = model.Gt("channel")
	}
	c.Param = strconv.FormatInt(o.ID, 10)
	c.URI = o.ReValue
	c.Badge = o.Corner
}
