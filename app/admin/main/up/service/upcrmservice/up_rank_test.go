package upcrmservice

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/up/model/upcrmmodel"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmservicesortRankInfo(t *testing.T) {
	convey.Convey("sortRankInfo", t, func(ctx convey.C) {
		var (
			planets  = []*upcrmmodel.UpRankInfo{}
			sortfunc sortRankFunc
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			sortRankInfo(planets, sortfunc)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestUpcrmservicesortByValueAsc(t *testing.T) {
	convey.Convey("sortByValueAsc", t, func(ctx convey.C) {
		var (
			p1 = &upcrmmodel.UpRankInfo{}
			p2 = &upcrmmodel.UpRankInfo{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := sortByValueAsc(p1, p2)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmservicesortByValueDesc(t *testing.T) {
	convey.Convey("sortByValueDesc", t, func(ctx convey.C) {
		var (
			p1 = &upcrmmodel.UpRankInfo{}
			p2 = &upcrmmodel.UpRankInfo{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := sortByValueDesc(p1, p2)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmservicerefreshUpRankDate(t *testing.T) {
	convey.Convey("refreshUpRankDate", t, func(ctx convey.C) {
		var (
			date = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.refreshUpRankDate(date)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestUpcrmserviceUpRankQueryList(t *testing.T) {
	convey.Convey("UpRankQueryList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &upcrmmodel.UpRankQueryArgs{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.UpRankQueryList(c, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}
