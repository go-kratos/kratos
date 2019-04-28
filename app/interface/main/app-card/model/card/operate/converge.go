package operate

import (
	"encoding/json"
	"strconv"

	"go-common/app/interface/main/app-card/model"
	"go-common/library/log"
)

type Converge struct {
	ID      int64           `json:"id,omitempty"`
	ReType  int             `json:"re_type,omitempty"`
	ReValue string          `json:"re_value,omitempty"`
	Title   string          `json:"title,omitempty"`
	Cover   string          `json:"cover,omitempty"`
	Content json.RawMessage `json:"content,omitempty"`
	// extra
	Items []*Converge `json:"contents,omitempty"`
	Goto  model.Gt    `json:"goto,omitempty"`
	Pid   int64       `json:"pid,omitempty"`
	Param string      `json:"param,omitempty"`
	URI   string      `json:"uri,omitempty"`
}

func (c *Converge) Change() {
	var contents []*struct {
		Ctype  string `json:"ctype,omitempty"`
		Cvalue string `json:"cvalue,omitempty"`
	}
	if err := json.Unmarshal(c.Content, &contents); err != nil {
		log.Error("%+v", err)
		return
	}
	c.Goto = model.OperateType[c.ReType]
	c.Param = c.ReValue
	c.URI = model.FillURI(c.Goto, c.Param, nil)
	c.Items = make([]*Converge, 0, len(contents))
	for _, content := range contents {
		var gt model.Gt
		id, _ := strconv.ParseInt(content.Cvalue, 10, 64)
		if id == 0 {
			continue
		}
		switch content.Ctype {
		case "0":
			gt = model.GotoAv
		case "1":
			gt = model.GotoLive
		case "2":
			gt = model.GotoArticle
		default:
			continue
		}
		c.Items = append(c.Items, &Converge{Pid: id, Goto: gt})
	}
}
