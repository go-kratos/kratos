package dao

import (
	"go-common/app/admin/main/tv/model"
	"testing"

	"time"

	"context"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSetModulePublishCache(t *testing.T) {
	var (
		c      = context.Background()
		pageID = "18"
		p      = model.ModPub{
			Time:  time.Now().Format("2006-01-02 15:04:05"),
			State: 1,
		}
	)
	convey.Convey("SetModPub", t, func(ctx convey.C) {
		err := d.SetModPub(c, pageID, p)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoGetModulePublishCache(t *testing.T) {
	var (
		c      = context.Background()
		pageID = "18"
	)
	convey.Convey("GetModPub", t, func(ctx convey.C) {
		p, err := d.GetModPub(c, pageID)
		ctx.Convey("Then err should be nil.p should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p, convey.ShouldNotBeNil)
		})
	})
}
