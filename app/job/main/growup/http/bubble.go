package http

import (
	"time"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func syncBubbleMetaTask(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	t, err := time.ParseInLocation("2006-01-02", v.Date, time.Local)
	if err != nil {
		log.Error("syncBubbleMetaTask date error(date:%s)", v.Date)
		return
	}
	err = svr.SyncIncomeBubbleMetaTask(c, t)
	if err != nil {
		log.Error("svr.SyncIncomeBubbleMetaTask error(%v)", err)
	}
	c.JSON(nil, err)
}

func syncBubbleMeta(c *bm.Context) {
	v := new(struct {
		Start string `form:"start" validate:"required"`
		End   string `form:"end" validate:"required"`
		BType int    `form:"b_type" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	start, err := time.ParseInLocation("2006-01-02", v.Start, time.Local)
	if err != nil {
		log.Error("syncBubbleMeta date error(date:%s)", v.Start)
		return
	}
	end, err := time.ParseInLocation("2006-01-02", v.End, time.Local)
	if err != nil {
		log.Error("syncBubbleMeta date error(date:%s)", v.End)
		return
	}
	rows, err := svr.SyncIncomeBubbleMeta(c, start, end, v.BType)
	if err != nil {
		log.Error("svr.SyncIncomeBubbleMeta arg(%v) error(%v)", v, err)
	}
	c.JSON(rows, err)
}

func snapshotBubbleTask(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	t, err := time.ParseInLocation("2006-01-02", v.Date, time.Local)
	if err != nil {
		log.Error("snapshotBubbleTask date error(date:%s)", v.Date)
		return
	}
	err = svr.SnapshotBubbleIncomeTask(c, t)
	if err != nil {
		log.Error("svr.SnapshotBubbleIncomeTask error(%v)", err)
	}
	c.JSON(nil, err)
}
