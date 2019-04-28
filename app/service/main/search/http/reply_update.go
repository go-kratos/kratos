package http

import (
	"encoding/json"

	"go-common/app/service/main/search/dao"
	"go-common/app/service/main/search/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// replyUpdate
func replyUpdate(c *bm.Context) {
	params := c.Request.Form
	appid := params.Get("appid")
	if appid == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	switch appid {
	case "reply_record":
		replyRecordUpdate(c)
	default:
		c.JSON(nil, ecode.RequestErr)
	}
}

func replyRecordUpdate(c *bm.Context) {
	var (
		err      error
		bulkItem []dao.BulkItem
		d        []*model.ReplyRecordUpdateParams
	)
	params := c.Request.Form
	data := params.Get("data")
	if data == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = json.Unmarshal([]byte(data), &d); err != nil {
		log.Error("json.Unmarshal error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, v := range d {
		bulkItem = append(bulkItem, v)
	}
	if err = svr.Update(c, "replyExternal", bulkItem); err != nil {
		log.Error("srv.Update error(%v)", err)
	}
	c.JSON(nil, err)
}
