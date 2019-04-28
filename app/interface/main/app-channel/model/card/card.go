package card

import (
	"encoding/json"
	"strconv"

	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/app/interface/main/app-channel/model"
	"go-common/app/interface/main/app-channel/model/recommend"
	"go-common/library/log"
)

type Card struct {
	ID         int64  `json:"-"`
	Title      string `json:"-"`
	ChannelID  int64  `json:"-"`
	Type       string `json:"-"`
	Value      int64  `json:"-"`
	Reason     string `json:"-"`
	ReasonType int8   `json:"-"`
	Pos        int    `json:"-"`
	FromType   string `json:"-"`
}

type CardPlat struct {
	CardID    int64  `json:"-"`
	Plat      int8   `json:"-"`
	Condition string `json:"-"`
	Build     int    `json:"-"`
}

type UpCard struct {
	ID         int64           `json:"id,omitempty"`
	Type       string          `json:"type,omitempty"`
	Title      string          `json:"title,omitempty"`
	RcmdReason string          `json:"rcmd_reason,omitempty"`
	Content    json.RawMessage `json:"content,omitempty"`
	Data       *recommend.Item `json:"data,omitempty"`
}

type UpContent struct {
	CType  string      `json:"ctype,omitempty"`
	CValue interface{} `json:"cvalue,omitempty"`
}

type ChannelSingle struct {
	AID       interface{} `json:"aid,omitempty"`
	ChannelID interface{} `json:"channel_id,omitempty"`
}

func (c *UpCard) Change() {
	data := []*UpContent{}
	if err := json.Unmarshal(c.Content, &data); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", c.Content, err)
		return
	}
	upItem := make([]*recommend.Item, 0, len(data))
	for _, d := range data {
		switch d.CType {
		case "mid", "channel_id":
			var (
				value int64
				err   error
			)
			switch v := d.CValue.(type) {
			case float64:
				value = int64(v)
			case string:
				if value, err = strconv.ParseInt(v, 10, 64); err != nil {
					continue
				}
			}
			item := &recommend.Item{ID: value}
			upItem = append(upItem, item)
		}
	}
	if len(upItem) < 3 {
		return
	}
	// TODO rcmd_reason
	c.Data = &recommend.Item{Goto: model.GotoSubscribe, ID: c.ID, Config: &recommend.Config{Title: c.RcmdReason}, Items: upItem}
}

func (c *UpCard) ChannelSingleChange() {
	data := &ChannelSingle{}
	if err := json.Unmarshal(c.Content, &data); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", c.Content, err)
		return
	}
	var (
		aid, channelid int64
		err            error
	)
	switch v := data.AID.(type) {
	case float64:
		aid = int64(v)
	case string:
		if aid, err = strconv.ParseInt(v, 10, 64); err != nil {
			return
		}
	}
	switch v := data.ChannelID.(type) {
	case float64:
		channelid = int64(v)
	case string:
		if channelid, err = strconv.ParseInt(v, 10, 64); err != nil {
			return
		}
	}
	c.Data = &recommend.Item{Goto: model.GotoChannelRcmd, ID: aid, Config: &recommend.Config{Title: c.Title}, TagID: channelid}
}

func (c *Card) CardToAiChange() (a *ai.Item) {
	a = &ai.Item{
		Goto:       c.Type,
		ID:         c.Value,
		RcmdReason: c.fromRcmdReason(),
	}
	return
}

func (c *Card) fromRcmdReason() (a *ai.RcmdReason) {
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
