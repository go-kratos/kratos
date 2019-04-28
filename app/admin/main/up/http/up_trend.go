package http

import (
	"go-common/app/admin/main/up/util"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"time"
)

func yesterday(c *bm.Context) {
	v := new(struct {
		Date int64 `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	res, err := Svc.Crmservice.QueryYesterday(c, time.Unix(v.Date, 0))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(res, err)
}

func trend(c *bm.Context) {
	v := new(struct {
		StatType int `form:"stat_type" validate:"required"`
		DateType int `form:"date_type" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
	}

	endDate := util.TruncateDate(time.Now())
	var days = 7
	if v.DateType == 2 {
		days = 30
	}

	data, err := Svc.Crmservice.QueryTrend(c, v.StatType, endDate, days)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(data, nil)
}

func trendDetail(c *bm.Context) {
	v := new(struct {
		DateType int `form:"date_type" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	endDate := util.TruncateDate(time.Now())
	var days = 7
	if v.DateType == 2 {
		days = 30
	}
	data, err := Svc.Crmservice.QueryDetail(c, endDate, days)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
