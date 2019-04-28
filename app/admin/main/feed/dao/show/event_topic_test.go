package show

import (
	"go-common/app/admin/main/feed/model/show"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestShowEventTopicAdd(t *testing.T) {
	convey.Convey("EventTopicAdd", t, func(ctx convey.C) {
		var (
			param = &show.EventTopicAP{
				Title:   "title",
				Desc:    "desc",
				Cover:   "cover",
				Retype:  1,
				Revalue: "revalue",
				Corner:  "corner",
				Person:  "person",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.EventTopicAdd(param)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestShowEventTopicUpdate(t *testing.T) {
	convey.Convey("EventTopicUpdate", t, func(ctx convey.C) {
		var (
			param = &show.EventTopicUP{
				ID:      1,
				Title:   "title111",
				Desc:    "desc111",
				Cover:   "cover111",
				Retype:  1,
				Revalue: "revalue111",
				Corner:  "corner111",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.EventTopicUpdate(param)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestShowEventTopicDelete(t *testing.T) {
	convey.Convey("EventTopicDelete", t, func(ctx convey.C) {
		var (
			id = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.EventTopicDelete(id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestShowETFindByID(t *testing.T) {
	convey.Convey("ETFindByID", t, func(ctx convey.C) {
		var (
			id = int64(2)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.ETFindByID(id)
			ctx.Convey("Then err should be nil.topic should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
