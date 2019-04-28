package show

import (
	"testing"

	"go-common/app/admin/main/feed/model/show"

	"github.com/smartystreets/goconvey/convey"
)

func TestShowSearchWebCardAdd(t *testing.T) {
	convey.Convey("SearchWebCardAdd", t, func(ctx convey.C) {
		var (
			param = &show.SearchWebCardAP{
				Type:    1,
				Title:   "搜索卡片",
				Desc:    "卡片描述",
				Cover:   "//bfs:",
				ReType:  1,
				ReValue: "http://",
				Corner:  "角标",
				Person:  "person",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SearchWebCardAdd(param)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestShowSearchWebCardUpdate(t *testing.T) {
	convey.Convey("SearchWebCardUpdate", t, func(ctx convey.C) {
		var (
			param = &show.SearchWebCardUP{
				ID:      1,
				Type:    1,
				Title:   "AA搜索卡片",
				Desc:    "AA卡片描述",
				Cover:   "//bfs:",
				ReType:  1,
				ReValue: "http://",
				Corner:  "角标",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SearchWebCardUpdate(param)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestShowSearchWebCardDelete(t *testing.T) {
	convey.Convey("SearchWebCardDelete", t, func(ctx convey.C) {
		var (
			id = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SearchWebCardDelete(id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestShowSWBFindByID(t *testing.T) {
	convey.Convey("SWBFindByID", t, func(ctx convey.C) {
		var (
			id = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.SWBFindByID(id)
			ctx.Convey("Then err should be nil.topic should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
