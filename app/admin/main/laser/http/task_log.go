package http

import (
	"go-common/app/admin/main/laser/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
)

func queryTaskLog(c *blademaster.Context) {
	v := new(struct {
		MID       int64  `form:"mid"`
		TaskID    int64  `form:"task_id"`
		Platform  int    `form:"platform"`
		TaskState int    `form:"task_state"`
		Sortby    string `form:"sort"`
		PageNo    int    `form:"page_no"`
		PageSize  int    `form:"page_size"`
	})
	err := c.Bind(v)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	logs, count, err := svc.QueryTaskLog(c, v.MID, v.TaskID, v.Platform, v.TaskState, v.Sortby, v.PageNo, v.PageSize)
	if err != nil {
		log.Error("svc.QueryTaskLog() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	pager := &model.TaskLogPager{
		PageNo:   v.PageNo,
		PageSize: v.PageSize,
		Items:    logs,
		Total:    count,
	}
	c.JSON(pager, nil)
}
