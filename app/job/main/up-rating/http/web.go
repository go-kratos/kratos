package http

import (
	"time"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func pastScore(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	date, err := time.Parse("2006-01-02", v.Date)
	if err != nil {
		log.Error("pastScore date error(date:%s)", v.Date)
		return
	}
	err = svr.RunPastScore(c, date)
	if err != nil {
		log.Error("svr.RunPastScore error(%v)", err)
	}
	c.JSON(nil, err)
}

func pastRecord(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.InsertPastRecord(c, v.Date)
	if err != nil {
		log.Error("svr.InsertPastRecord error(%v)", err)
	}
	c.JSON(nil, err)
}

func score(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	date, err := time.Parse("2006-01-02", v.Date)
	if err != nil {
		log.Error("score date error(date:%s)", v.Date)
		return
	}
	err = svr.Run(c, date)
	if err != nil {
		log.Error("svr.Run error(%v)", err)
	}
	c.JSON(nil, err)
}

func delTrends(c *bm.Context) {
	v := new(struct {
		Table string `form:"table" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.DelTrends(c, v.Table)
	if err != nil {
		log.Error("svr.DelTrends error(%v)", err)
	}
	c.JSON(nil, err)
}

func delScore(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	date, err := time.Parse("2006-01-02", v.Date)
	if err != nil {
		log.Error("del score date error(date:%s)", v.Date)
		return
	}
	err = svr.DelRatings(c, date)
	if err != nil {
		log.Error("svr.DelRatings error(%v)", err)
	}
	c.JSON(nil, err)
}

func statistics(c *bm.Context) {
	var err error
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	date, err := time.Parse("2006-01-02", v.Date)
	if err != nil {
		log.Error("statistics date error(date:%s)", v.Date)
		return
	}
	err = svr.RunStatistics(c, date)
	if err != nil {
		log.Error("svr.RunStatistics error(%v)", err)
	}
	c.JSON(nil, err)
}

func trend(c *bm.Context) {
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	date, err := time.Parse("2006-01-02", v.Date)
	if err != nil {
		log.Error("trend date error(date:%s)", v.Date)
		return
	}
	err = svr.CalTrend(c, date)
	if err != nil {
		log.Error("svr.Trend error(%v)", err)
	}
	c.JSON(nil, err)
}

func taskStatus(c *bm.Context) {
	var err error
	v := new(struct {
		Date    string `form:"date" validate:"required"`
		Type    int    `form:"type"`
		Status  int    `form:"status"`
		Message string `form:"message"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if err = svr.InsertTaskStatus(c, v.Type, v.Status, v.Date, v.Message); err != nil {
		log.Error("svr.UpdateTaskStatus error(%v)", err)
	}
	c.JSON(nil, err)
}
