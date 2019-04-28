package show

import (
	"go-common/app/admin/main/feed/model/show"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestShowPopularStarsAdd(t *testing.T) {
	convey.Convey("PopularStarsAdd", t, func(ctx convey.C) {
		var (
			param = &show.PopularStarsAP{
				Type:      "up_rcmd_new",
				Value:     "31",
				Title:     "title",
				LongTitle: "longtitle",
				Content:   "[{\"id\":10099666,\"title\":\"uat-case-新转码系统1515660475\"},{\"id\":10099667,\"title\":\"测试-test22-秒开-1515660506\"}]",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.PopularStarsAdd(param)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestShowPopularStarsUpdate(t *testing.T) {
	convey.Convey("PopularStarsUpdate", t, func(ctx convey.C) {
		var (
			param = &show.PopularStarsUP{
				ID:        16,
				Type:      "up_rcmd_new",
				Value:     "32",
				Title:     "title",
				LongTitle: "longtitle",
				Content:   "[{\"id\":10099666,\"title\":\"uat-case-新转码系统1515660475\"},{\"id\":10099667,\"title\":\"测试-test22-秒开-1515660506\"}]",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.PopularStarsUpdate(param)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestShowPopularStarsDelete(t *testing.T) {
	convey.Convey("PopularStarsDelete", t, func(ctx convey.C) {
		var (
			id = int64(11)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.PopularStarsDelete(id, "up_rcmd_new")
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
