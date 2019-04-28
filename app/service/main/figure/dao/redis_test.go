package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/service/main/figure/model"
)

func TestDaofigureKey(t *testing.T) {
	convey.Convey("figureKey", t, func(ctx convey.C) {
		var (
			mid = int64(46333)
		)
		ctx.Convey("When everything right.", func(ctx convey.C) {
			key := figureKey(mid)
			ctx.Convey("Then key should equal f:key.", func(ctx convey.C) {
				ctx.So(key, convey.ShouldEqual, "f:46333")
			})
		})
	})
}

func TestDaoPingRedis(t *testing.T) {
	convey.Convey("PingRedis", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything right.", func(ctx convey.C) {
			err := d.PingRedis(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddFigureInfoCache(t *testing.T) {
	convey.Convey("AddFigureInfoCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			figure = &model.Figure{
				Mid:             46333,
				Score:           2333,
				LawfulScore:     123,
				WideScore:       321,
				FriendlyScore:   19999,
				BountyScore:     1,
				CreativityScore: 0,
				Ver:             2333,
			}
		)
		ctx.Convey("When add FigureInfoCache.", func(ctx convey.C) {
			err := d.AddFigureInfoCache(c, figure)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Convey("When get FigureInfoCache.", func(ctx convey.C) {
					figure2, err := d.FigureInfoCache(c, figure.Mid)
					ctx.Convey("Then err should be nil.figure2 should resemble figure.", func(ctx convey.C) {
						ctx.So(err, convey.ShouldBeNil)
						ctx.So(figure2, convey.ShouldResemble, figure)
					})
				})
				ctx.Convey("When get FigureBatchInfoCache.", func(ctx convey.C) {
					figures, missIndex, err := d.FigureBatchInfoCache(c, []int64{figure.Mid})
					ctx.Convey("Then err should be nil.missIndex should be empty.figures should have length 1.figuers[0] should resemble figure", func(ctx convey.C) {
						ctx.So(err, convey.ShouldBeNil)
						ctx.So(missIndex, convey.ShouldBeEmpty)
						ctx.So(figures, convey.ShouldHaveLength, 1)
						ctx.So(figures[0], convey.ShouldResemble, figure)
					})
				})
			})
		})
	})
}
