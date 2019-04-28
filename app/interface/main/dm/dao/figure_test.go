package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaofigureInfoURI(t *testing.T) {
	convey.Convey("figureInfoURI", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := testDao.figureInfoURI()
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoFigureInfo(t *testing.T) {
	convey.Convey("FigureInfo", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			score, err := testDao.FigureInfo(c, mid)
			convCtx.Convey("Then err should be nil.score should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(score, convey.ShouldNotBeNil)
			})
		})
	})
}
