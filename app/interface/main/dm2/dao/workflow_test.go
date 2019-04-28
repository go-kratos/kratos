package dao

import (
	"context"
	"go-common/app/interface/main/dm2/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoWorkFlowTagList(t *testing.T) {
	convey.Convey("WorkFlowTagList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			bid = int64(0)
			rid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.WorkFlowTagList(c, bid, rid)
		})
	})
}

func TestDaoWorkFlowAppealAdd(t *testing.T) {
	convey.Convey("WorkFlowAppealAdd", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &model.WorkFlowAppealAddReq{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.WorkFlowAppealAdd(c, req)
		})
	})
}

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
