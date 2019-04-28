package upcrmdao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"time"
)

func TestUpcrmdaoAsPrScore(t *testing.T) {
	convey.Convey("AsPrScore", t, func(ctx convey.C) {
		var info *UpQualityInfo
		history := info.AsPrScore()
		ctx.Convey("Then history should not be nil.", func(ctx convey.C) {
			ctx.So(history, convey.ShouldNotBeNil)
		})
	})
}

func TestUpcrmdaoAsQualityScore(t *testing.T) {
	convey.Convey("AsQualityScore", t, func(ctx convey.C) {
		var info *UpQualityInfo
		history := info.AsQualityScore()
		ctx.Convey("Then history should not be nil.", func(ctx convey.C) {
			ctx.So(history, convey.ShouldNotBeNil)
		})
	})
}

func TestUpcrmdaoUpdateCreditScore(t *testing.T) {
	var (
		score = int(0)
		mid   = int64(0)
	)
	convey.Convey("UpdateCreditScore", t, func(ctx convey.C) {
		affectRow, err := d.UpdateCreditScore(score, mid)
		err = IgnoreErr(err)
		ctx.Convey("Then err should be nil.affectRow should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affectRow, convey.ShouldNotBeNil)
		})
	})
}

func TestUpcrmdaoUpdateQualityAndPrScore(t *testing.T) {
	var (
		prScore      = int(0)
		qualityScore = int(0)
		mid          = int64(0)
	)
	convey.Convey("UpdateQualityAndPrScore", t, func(ctx convey.C) {
		affectRow, err := d.UpdateQualityAndPrScore(prScore, qualityScore, mid)
		err = IgnoreErr(err)
		ctx.Convey("Then err should be nil.affectRow should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affectRow, convey.ShouldNotBeNil)
		})
	})
}

func TestUpcrmdaoInsertScoreHistory(t *testing.T) {
	var (
		info = &UpQualityInfo{}
	)
	convey.Convey("InsertScoreHistory", t, func(ctx convey.C) {
		affectRow, err := d.InsertScoreHistory(info)
		err = IgnoreErr(err)
		ctx.Convey("Then err should be nil.affectRow should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affectRow, convey.ShouldNotBeNil)
		})
	})
}

func TestUpcrmdaoInsertBatchScoreHistory(t *testing.T) {
	var (
		infoList = []*UpQualityInfo{{Mid: 100, Cdate: time.Now().Format(TimeFmtDate)}}
		tablenum = int(0)
	)
	convey.Convey("InsertBatchScoreHistory", t, func(ctx convey.C) {
		affectRow, err := d.InsertBatchScoreHistory(infoList, tablenum)
		err = IgnoreErr(err)
		ctx.Convey("Then err should be nil.affectRow should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affectRow, convey.ShouldNotBeNil)
		})
	})
}
