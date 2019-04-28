package http

import (
	"time"

	"go-common/library/log"

	bm "go-common/library/net/http/blademaster"
)

func combineMails(c *bm.Context) {
	log.Info("begin combineMails")
	v := new(struct {
		Date string `form:"date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	t, err := time.Parse("2006-01-02", v.Date)
	if err != nil {
		log.Error("sendIncome date error!date:%s", v.Date)
		return
	}
	err = svr.CombineMailsByHTTP(c, t.Year(), int(t.Month()), t.Day())
	if err != nil {
		log.Error("Exec CombineMailsByHTTP error!(%v)", err)
	} else {
		log.Info("Exec CombineMailsByHTTP succeed!")
	}
	c.JSON(nil, err)
}
