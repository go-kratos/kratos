package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/tv/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSetSearInterRankCache(t *testing.T) {
	var (
		c    = context.Background()
		rank = []*model.OutSearchInter{}
	)
	convey.Convey("SetSearchInterv", t, func(ctx convey.C) {
		rank = append(rank, &model.OutSearchInter{
			Keyword: "key",
			Status:  "1",
		})
		err := d.SetSearchInterv(c, rank)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoGetSearInterRankCache(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("GetSearchInterv", t, func(ctx convey.C) {
		rank, err := d.GetSearchInterv(c)
		ctx.Convey("Then err should be nil.rank should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rank, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetPublishCache(t *testing.T) {
	var (
		c     = context.Background()
		state = &model.PublishStatus{
			Time:  time.Now().Format("2006-01-02 15:04:05"),
			State: 1,
		}
	)
	convey.Convey("SetPublishCache", t, func(ctx convey.C) {
		err := d.SetPublishCache(c, state)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoGetPublishCache(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("GetPublishCache", t, func(ctx convey.C) {
		state, err := d.GetPublishCache(c)
		ctx.Convey("Then err should be nil.state should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(state, convey.ShouldNotBeNil)
		})
	})
}
