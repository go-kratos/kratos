package show

import (
	"testing"

	"go-common/app/admin/main/feed/model/show"

	"github.com/smartystreets/goconvey/convey"
)

func TestShowSearchWebAdd(t *testing.T) {
	convey.Convey("SearchWebAdd", t, func(ctx convey.C) {
		var (
			param = &show.SearchWebAP{
				CardType:    1,
				CardValue:   "10",
				Stime:       1545701985,
				Etime:       1545711985,
				Priority:    1,
				Person:      "quguolin",
				ApplyReason: "test",
				Query:       "[{\"id\":7,\"value\":\"test1\"},{\"id\":8,\"value\":\"test2\"}]",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SearchWebAdd(param)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestShowSearchWebUpdate(t *testing.T) {
	convey.Convey("SearchWebUpdate", t, func(ctx convey.C) {
		var (
			param = &show.SearchWebUP{
				ID:          2,
				CardType:    1,
				CardValue:   "10",
				Stime:       1545701985,
				Etime:       1545711985,
				Check:       1,
				Status:      1,
				Priority:    1,
				Person:      "quguolin",
				ApplyReason: "test",
				Query:       "[{\"id\":7,\"value\":\"test1\"},{\"id\":8,\"value\":\"test2\"},{\"sid\":10099668,\"value\":\"aaa\"}]",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SearchWebUpdate(param)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestShowSearchWebDelete(t *testing.T) {
	convey.Convey("SearchWebDelete", t, func(ctx convey.C) {
		var (
			id = int64(4)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SearchWebDelete(id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestShowSearchWebOption(t *testing.T) {
	convey.Convey("SearchWebOption", t, func(ctx convey.C) {
		var (
			up = &show.SearchWebOption{
				ID:     1,
				Check:  4,
				Status: 1,
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SearchWebOption(up)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestShowSWTimeValid(t *testing.T) {
	convey.Convey("SWTimeValid", t, func(ctx convey.C) {
		var (
			up = &show.SWTimeValid{
				Priority: 1,
				Query:    "test1",
				STime:    1543190400,
				ETime:    1543449600,
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.SWTimeValid(up)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestShowSWFindByID(t *testing.T) {
	convey.Convey("SWFindByID", t, func(ctx convey.C) {
		var (
			id = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			d.SWFindByID(id)
		})
	})
}
