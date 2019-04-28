package http

import (
	"encoding/json"

	"go-common/app/admin/main/search/dao"
	"go-common/app/admin/main/search/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// updateBlocked .
func updateArchive(c *bm.Context) {
	form := c.Request.Form
	appid := form.Get("appid")

	switch appid {
	case "task_qa_fans":
		updateTaskQaFans(c)
	default:
		c.JSON(nil, ecode.RequestErr)
	}
}

func updateTaskQaFans(c *bm.Context) {
	var (
		err      error
		bulkItem []dao.BulkItem
		d        []*model.TaskQaFansParams
		form     = c.Request.Form
	)
	data := form.Get("data")
	if data == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = json.Unmarshal([]byte(data), &d); err != nil {
		log.Error("json.Unmarshal error(%v)", err)
		c.JSON(nil, err)
		return
	}
	for _, v := range d {
		bulkItem = append(bulkItem, v)
	}
	if err = svr.Update(c, "ssd_archive", bulkItem); err != nil {
		log.Error("srv.Update error(%v)", err)
	}
	c.JSON(nil, err)
}
