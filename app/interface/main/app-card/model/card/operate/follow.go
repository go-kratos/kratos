package operate

import (
	"encoding/json"

	"go-common/app/interface/main/app-card/model"
	"go-common/library/log"
)

type Follow struct {
	ID      int64           `json:"id,omitempty"`
	Type    string          `json:"type,omitempty"`
	Title   string          `json:"title,omitempty"`
	Content json.RawMessage `json:"content,omitempty"`
	// extra
	Items []*Follow `json:"items,omitempty"`
	Goto  model.Gt  `json:"goto,omitempty"`
	Pid   int64     `json:"pid,omitempty"`
	Tid   int64     `json:"tid,omitempty"`
}

func (c *Follow) Change() {
	switch c.Type {
	case "upper", "channel_three":
		var contents []*struct {
			Ctype  string `json:"ctype,omitempty"`
			Cvalue int64  `json:"cvalue,omitempty"`
		}
		if err := json.Unmarshal(c.Content, &contents); err != nil {
			log.Error("%+v", err)
			return
		}
		items := make([]*Follow, 0, len(contents))
		for _, content := range contents {
			item := &Follow{Type: content.Ctype, Pid: content.Cvalue}
			switch content.Ctype {
			case "mid":
				item.Goto = model.GotoMid
			case "channel_id":
				item.Goto = model.GotoTag
			}
			items = append(items, item)
		}
		if len(items) < 3 {
			return
		}
		c.Items = items
	case "channel_single":
		var content struct {
			Aid       int64 `json:"aid"`
			ChannelID int64 `json:"channel_id"`
		}
		if err := json.Unmarshal(c.Content, &content); err != nil {
			log.Error("%+v", err)
			return
		}
		c.Pid = content.Aid
		c.Tid = content.ChannelID
	}
}
