package http

import (
	"fmt"

	"go-common/app/service/main/up/service"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func baseStat(c *bm.Context) {
	var arg = new(struct {
		Mid  int64  `form:"mid" validate:"required"`
		Date string `form:"date" validate:"required"`
	})
	var res interface{}
	var err error
	var errMsg string
	switch {
	default:
		if err = c.Bind(arg); err != nil {
			log.Error("request argument bind fail, err=%v", err)
			errMsg = fmt.Sprintf("wrong argument, %s", err.Error())
			err = ecode.RequestErr
			break
		}

		var date = arg.Date
		var mid = arg.Mid
		var d, e = Svc.Data.BaseUpStat(c, mid, date)
		err = e
		if err != nil {
			log.Error("get hbase fail, mid=%d, err=%v", mid, err)
			return
		}
		log.Info("get from hbase ok, mid=%d, stat=%+v", mid, d)
		res = map[string]interface{}{
			"stat": d,
		}
	}

	if err != nil {
		service.BmHTTPErrorWithMsg(c, err, errMsg)
	} else {
		c.JSON(res, err)
	}
}
