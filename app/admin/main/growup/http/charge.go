package http

import (
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func archiveChargeStatis(c *bm.Context) {
	v := new(struct {
		CategoryID []int64 `form:"category_id,split"`
		Type       int     `form:"type"`
		GroupType  int     `form:"group_type" default:"1"`
		FromTime   int64   `form:"from_time" validate:"required,min=1"`
		ToTime     int64   `form:"to_time" validate:"required,min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	data, err := incomeSvr.ArchiveChargeStatis(c, v.CategoryID, v.Type, v.GroupType, v.FromTime, v.ToTime)
	if err != nil {
		log.Error("incomeSvr.ArchiveChargeStatis error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func archiveChargeSection(c *bm.Context) {
	v := new(struct {
		CategoryID []int64 `form:"category_id,split"`
		Type       int     `form:"type"`
		GroupType  int     `form:"group_type" default:"1"`
		FromTime   int64   `form:"from_time" validate:"required,min=1"`
		ToTime     int64   `form:"to_time" validate:"required,min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	data, err := incomeSvr.ArchiveChargeSection(c, v.CategoryID, v.Type, v.GroupType, v.FromTime, v.ToTime)
	if err != nil {
		log.Error("growup incomeSvr.ArchiveChargeSection error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func archiveChargeDetail(c *bm.Context) {
	v := new(struct {
		AID  int64 `form:"aid" validate:"required"`
		Type int   `form:"type"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	data, err := incomeSvr.ArchiveChargeDetail(c, v.AID, v.Type)
	if err != nil {
		log.Error("growup incomeSvr.ArchiveDetail error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func bgmChargeDetail(c *bm.Context) {
	v := new(struct {
		SID int64 `form:"sid" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	data, err := incomeSvr.BgmChargeDetail(c, v.SID)
	if err != nil {
		log.Error("growup incomeSvr.BgmChargeDetail error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func upRatio(c *bm.Context) {
	v := new(struct {
		From  int64 `form:"from"`
		Limit int64 `form:"limit" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	data, err := incomeSvr.UpRatio(c, v.From, v.Limit)
	if err != nil {
		log.Error("growup incomeSvr.UpRatio error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
