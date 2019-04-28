package http

import (
	"time"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func tagIncome(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := tagSvr.TagIncomeAll(c, v.Date)
	if err != nil {
		log.Error("tagSvr.TagIncomeAll error(%v)", err)
	}
	c.JSON(nil, err)
}

func tagRatio(c *bm.Context) {
	log.Info("begin update tag ratio")
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})

	if err := c.Bind(v); err != nil {
		return
	}
	err := tagSvr.TagRatioAll(c, v.Date)
	if err != nil {
		log.Error("svr.TagRatio error(%v)", err)
	} else {
		log.Info("tag ratio succeed!")
	}
	c.JSON(nil, err)
}

func tagUps(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})

	if err := c.Bind(v); err != nil {
		return
	}
	date, err := time.ParseInLocation("2006-01-02", v.Date, time.Local)
	if err != nil {
		log.Error("time.Parse date(%s) error", v.Date)
		c.JSON(nil, err)
		return
	}

	err = tagSvr.TagUps(c, date)
	if err != nil {
		log.Error("svr.TagUps error(%v)", err)
	} else {
		log.Info("tag ups succeed!")
	}
	c.JSON(nil, err)
}

func tagExtraIncome(c *bm.Context) {
	v := new(struct {
		TagIDs []int64 `form:"tag_ids,split" validate:"required"`
		Tag    string  `form:"tag" validate:"required"`
		Ratio  int64   `form:"ratio"`
		Start  string  `form:"start" validate:"required"`
		End    string  `form:"end" validate:"required"`
	})

	if err := c.Bind(v); err != nil {
		return
	}
	err := tagSvr.TagExtraIncome(c, v.TagIDs, v.Tag, v.Ratio, v.Start, v.End)
	if err != nil {
		log.Error("tagSvr.TagExtraIncome error(%v)", err)
	}
	c.JSON(nil, err)
}
