package converge

import (
	"encoding/json"
	"go-common/library/log"
	"strconv"

	"go-common/app/interface/main/app-channel/model"
)

type Card struct {
	ID       int64
	ReType   int
	ReValue  string
	Title    string
	Cover    string
	Content  json.RawMessage
	Contents []*Content
}

type JSONContent struct {
	CType  string          `json:"ctype"`
	CTitle string          `json:"ctitle"`
	CValue json.RawMessage `json:"cvalue"`
}

type Content struct {
	Goto  string `json:"goto,omitempty"`
	ID    int64  `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

func (c *Card) Change() {
	var (
		contents     []*Content
		jsonContents []*JSONContent
		err          error
	)
	if err = json.Unmarshal(c.Content, &jsonContents); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", c.Content, err)
		return
	}
	contents = make([]*Content, 0, len(jsonContents))
	for _, js := range jsonContents {
		switch js.CType {
		case "0", "1", "2":
			content := &Content{Title: js.CTitle}
			if js.CType == "0" {
				content.Goto = model.GotoAv
			} else if js.CType == "1" {
				content.Goto = model.GotoLive
			} else if js.CType == "2" {
				content.Goto = model.GotoArticle
			} else {
				continue
			}
			var (
				idStr string
				id    int64
			)
			if err = json.Unmarshal(js.CValue, &idStr); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", js.CValue, err)
				continue
			}
			if id, _ = strconv.ParseInt(idStr, 10, 64); id != 0 {
				content.ID = id
				contents = append(contents, content)
			}
		}
	}
	c.Contents = contents
}
