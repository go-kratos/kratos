package http

import (
	"time"

	"go-common/app/admin/main/app/model"
	"go-common/app/admin/main/app/model/notice"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// notices select notice all
func notices(c *bm.Context) {
	c.JSON(noticeSvc.Notices(c))
}

// noticeByID select notice by id
func noticeByID(c *bm.Context) {
	v := &notice.Param{}
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(noticeSvc.NoticeByID(c, v.ID))
}

// addOrupdate add or update notice
func addOrupdate(c *bm.Context) {
	var (
		err error
		v   = &notice.Param{}
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Eftime, err = model.NoticeChangeTime(v.EftimeStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Extime, err = model.NoticeChangeTime(v.ExtimeStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.ID > 0 {
		err = noticeSvc.UpdateNotice(c, v, time.Now())
	} else {
		err = noticeSvc.Insert(c, v, time.Now())
	}
	c.JSON(nil, err)
}

// updateBuild modify build and conditions
func updateBuild(c *bm.Context) {
	v := &notice.Param{}
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, noticeSvc.UpdateBuild(c, v, time.Now()))
}

// updateState modify notice state
func updateState(c *bm.Context) {
	v := &notice.Param{}
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, noticeSvc.UpdateState(c, v, time.Now()))
}
