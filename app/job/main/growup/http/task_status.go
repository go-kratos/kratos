package http

import (
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func updateTaskStatus(c *bm.Context) {
	v := new(struct {
		Status int    `form:"status" validate:"required"`
		Date   string `form:"date" validate:"required"`
		Type   int    `form:"type"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.UpdateTaskStatus(c, v.Date, v.Type, v.Status)
	if err != nil {
		log.Error("svr.UpdateTaskStatus error(%v)", err)
	}
	c.JSON(nil, err)
}

func checkTaskColumn(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := chargeSrv.CheckTaskColumn(c, v.Date)
	if err != nil {
		log.Error("incomeSrv.CheckTaskColumn error(%v)", err)
	}
	c.JSON(nil, err)
}
