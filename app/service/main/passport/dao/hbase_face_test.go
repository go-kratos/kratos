package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/tsuna/gohbase/hrpc"
)

func TestDaoFaceApplies(t *testing.T) {
	var (
		c        = context.TODO()
		mid      = int64(100)
		from, to = int64(0), int64(_uint32Max)
		status   = ""
		operator = ""
	)
	convey.Convey("FaceApplies", t, func(ctx convey.C) {
		_, err := d.FaceApplies(c, mid, from, to, status, operator)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			//ctx.So(res, convey.ShouldNotBeNil)
			//
			//lastTs := int64(0)
			//lastID := int64(0)
			//for _, v := range res {
			//	ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			//		mt, err := time.ParseInLocation(_inputFormat, v.ModifyTime, _loc)
			//		ctx.So(err, convey.ShouldBeNil)
			//
			//		id, err := strconv.ParseInt(v.Operator, 10, 64)
			//		ctx.So(err, convey.ShouldBeNil)
			//
			//		// check ts seq
			//		ts := mt.Unix()
			//		if lastTs > 0 {
			//			ctx.So(lastTs, convey.ShouldBeGreaterThanOrEqualTo, ts)
			//		}
			//
			//		// check id seq when ts equal
			//		if lastTs == ts && lastID > 0 {
			//			ctx.So(lastID, convey.ShouldBeGreaterThan, id)
			//		}
			//
			//		lastTs = ts
			//		lastID = id
			//	})
			//}

		})
	})
}

func TestDaorowKeyFaceApplyMts(t *testing.T) {
	var (
		midStr = "123"
		mts    = int64(0)
		id     = int64(0)
	)
	convey.Convey("rowKeyFaceApplyMts", t, func(ctx convey.C) {
		p1 := rowKeyFaceApplyMts(midStr, mts, id)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoscanFaceRecord(t *testing.T) {
	var (
		cells = []*hrpc.Cell{
			{Family: []byte("c"), Qualifier: []byte("mid"), Value: []byte("123")},
			{Family: []byte("c"), Qualifier: []byte("at"), Value: []byte("1530846373")},
			{Family: []byte("c"), Qualifier: []byte("mt"), Value: []byte("1530846373")},
			{Family: []byte("c"), Qualifier: []byte("nf"), Value: []byte("nf")},
			{Family: []byte("c"), Qualifier: []byte("of"), Value: []byte("of")},
			{Family: []byte("c"), Qualifier: []byte("op"), Value: []byte("op")},
			{Family: []byte("c"), Qualifier: []byte("s"), Value: []byte("0")},
		}
	)
	convey.Convey("scanFaceRecord", t, func(ctx convey.C) {
		res, err := scanFaceRecord(cells)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
