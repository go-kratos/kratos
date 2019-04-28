package upcrmservice

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/up/model/upcrmmodel"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmserviceconvertValueList(t *testing.T) {
	convey.Convey("convertValueList", t, func(ctx convey.C) {
		var (
			score upcrmmodel.ScoreSectionHistory
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := convertValueList(score)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmservicegenerateScoreQueryXAxis(t *testing.T) {
	convey.Convey("generateScoreQueryXAxis", t, func(ctx convey.C) {
		var (
			num = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			axis := generateScoreQueryXAxis(num)
			ctx.Convey("Then axis should not be nil.", func(ctx convey.C) {
				ctx.So(axis, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceScoreQuery(t *testing.T) {
	convey.Convey("ScoreQuery", t, func(ctx convey.C) {
		var (
			context = context.Background()
			arg     = &upcrmmodel.ScoreQueryArgs{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.ScoreQuery(context, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmservicecalcScoreInfo(t *testing.T) {
	convey.Convey("calcScoreInfo", t, func(ctx convey.C) {
		var (
			datamap  map[int8]map[time.Time]upcrmmodel.UpScoreHistory
			stype    = int8(0)
			todate   = time.Now()
			fromdate = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			info := calcScoreInfo(datamap, stype, todate, fromdate)
			ctx.Convey("Then info should not be nil.", func(ctx convey.C) {
				ctx.So(info, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmservicegenerateDataMap(t *testing.T) {
	convey.Convey("generateDataMap", t, func(ctx convey.C) {
		var (
			scoreHistory = []upcrmmodel.UpScoreHistory{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := generateDataMap(scoreHistory)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmservicegetDataFromMap(t *testing.T) {
	convey.Convey("getDataFromMap", t, func(ctx convey.C) {
		var (
			dataMap   map[int8]map[time.Time]upcrmmodel.UpScoreHistory
			scoreType = int(0)
			date      = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, ok := getDataFromMap(dataMap, scoreType, date)
			ctx.Convey("Then data,ok should not be nil.", func(ctx convey.C) {
				ctx.So(ok, convey.ShouldNotBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceScoreQueryUp(t *testing.T) {
	convey.Convey("ScoreQueryUp", t, func(ctx convey.C) {
		var (
			context = context.Background()
			arg     = &upcrmmodel.ScoreQueryUpArgs{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.ScoreQueryUp(context, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceScoreQueryUpHistory(t *testing.T) {
	convey.Convey("ScoreQueryUpHistory", t, func(ctx convey.C) {
		var (
			context = context.Background()
			arg     = &upcrmmodel.ScoreQueryUpHistoryArgs{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.ScoreQueryUpHistory(context, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}
