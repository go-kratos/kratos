package http

import (
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func archiveIncome(c *bm.Context) {
	v := new(struct {
		Type int `form:"type"`
		Page int `form:"page" default:"1" validate:"min=1"`
		Size int `form:"size" default:"20" validate:"min=1"`
		All  int `form:"all" default:"0"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)

	data, err := svc.ArchiveIncome(c, mid, v.Type, v.Page, v.Size, v.All)
	if err != nil {
		log.Error("growup svc.ArchiveIncome error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func upSummary(c *bm.Context) {
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)

	data, err := svc.UpSummary(c, mid)
	if err != nil {
		log.Error("growup svc.UpSummary error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func archiveSummary(c *bm.Context) {
	v := new(struct {
		Type int `form:"type"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)

	data, err := svc.ArchiveSummary(c, v.Type, mid)
	if err != nil {
		log.Error("growup svc.ArchiveSummary error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func archiveDetail(c *bm.Context) {
	v := new(struct {
		Type      int   `form:"type"`
		ArchiveID int64 `form:"archive_id" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	data, err := svc.ArchiveDetail(c, v.Type, v.ArchiveID)
	if err != nil {
		log.Error("growup svc.ArchiveDetail error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func archiveBreach(c *bm.Context) {
	v := new(struct {
		Type int `form:"type"`
		Page int `form:"page" default:"1" validate:"min=1"`
		Size int `form:"size" default:"20" validate:"min=1"`
		All  int `form:"all" default:"0"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)

	data, err := svc.ArchiveBreach(c, mid, v.Type, v.Page, v.Size, v.All)
	if err != nil {
		log.Error("growup svc.ArchiveBreach error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func upIncomeStat(c *bm.Context) {
	v := new(struct {
		Type int `form:"type"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	data, err := svc.UpIncomeStat(c, v.Type, mid, time.Now())
	if err != nil {
		log.Error("growup svc.UpIncomeStat error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
