package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoWorkFlowAppealDelete(t *testing.T) {
	convey.Convey("WorkFlowAppealDelete", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			bid        = int64(0)
			oid        = int64(0)
			subtitleID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.WorkFlowAppealDelete(c, bid, oid, subtitleID)
		})
	})
}
