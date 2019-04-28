package search

import (
	"context"
	"testing"

	mdlSearch "go-common/app/interface/main/tv/model/search"

	"github.com/smartystreets/goconvey/convey"
)

func TestSearchUserSearch(t *testing.T) {
	convey.Convey("UserSearch", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &mdlSearch.UserSearch{
				Keyword:    "lex",
				Build:      "111",
				SearchType: "bili_user",
				Page:       1,
				Pagesize:   20,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			user, err := d.UserSearch(c, arg)
			ctx.Convey("Then err should be nil.user should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(user, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestSearchSearchAllWild(t *testing.T) {
	convey.Convey("SearchAllWild", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &mdlSearch.UserSearch{
				Keyword:    "工作细胞",
				Build:      "111",
				SearchType: "all",
				Page:       1,
				Pagesize:   20,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			user, err := d.SearchAllWild(c, arg)
			ctx.Convey("Then err should be nil.user should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(user, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestSearchPgcSearch(t *testing.T) {
	convey.Convey("PgcSearch", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &mdlSearch.UserSearch{
				Keyword:    "白兔糖",
				Build:      "111",
				SearchType: "all",
				Page:       1,
				Pagesize:   20,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			st, err := d.PgcSearch(c, arg)
			ctx.Convey("Then err should be nil.st should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(st, convey.ShouldNotBeNil)
			})
		})
	})
}
