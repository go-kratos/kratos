package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceAddAuditLog(t *testing.T) {
	convey.Convey("AddAuditLog", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			bizID   = int(0)
			tp      = int8(0)
			action  = ""
			uid     = int64(0)
			uname   = ""
			oids    = []int64{}
			index   = []interface{}{}
			content map[string]interface{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := svr.AddAuditLog(c, bizID, tp, action, uid, uname, oids, index, content)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
