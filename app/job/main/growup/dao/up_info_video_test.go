package dao

import (
	"context"
	//	"go-common/library/time"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetUpStateByMID(t *testing.T) {
	convey.Convey("GetUpStateByMID", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			state, err := d.GetUpStateByMID(c, mid)
			ctx.Convey("Then err should be nil.state should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(state, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUpCreditScore(t *testing.T) {
	convey.Convey("GetUpCreditScore", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{100}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			scores, err := d.GetUpCreditScore(c, mids)
			ctx.Convey("Then err should be nil.scores should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(scores, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpInfoVideo(t *testing.T) {
	convey.Convey("UpInfoVideo", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			offset = int64(0)
			limit  = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			last, ups, err := d.UpInfoVideo(c, offset, limit)
			ctx.Convey("Then err should be nil.last,ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ups, convey.ShouldNotBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMIDsByState(t *testing.T) {
	convey.Convey("MIDsByState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			state = int(100)
			table = "up_info_video"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.MIDsByState(c, state, table)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMIDsByStateType(t *testing.T) {
	convey.Convey("MIDsByStateType", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			typ   = int(2)
			state = int(3)
			table = "up_info_video"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.MIDsByStateType(c, typ, state, table)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateAccountState(t *testing.T) {
	convey.Convey("UpdateAccountState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			state = int(2)
			mids  = []int64{100}
			table = "up_info_video"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpdateAccountState(c, state, mids, table)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetDateSignedUps(t *testing.T) {
	convey.Convey("GetDateSignedUps", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			startAt = time.Now()
			endAt   = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.GetDateSignedUps(c, startAt, endAt)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAllSignedUps(t *testing.T) {
	convey.Convey("GetAllSignedUps", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			data = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.GetAllSignedUps(c, data)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetVideoApplyUpCount(t *testing.T) {
	convey.Convey("GetVideoApplyUpCount", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			startAt = time.Now()
			endAt   = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.GetVideoApplyUpCount(c, startAt, endAt)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUpBaseInfo(t *testing.T) {
	convey.Convey("GetUpBaseInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = []int64{100}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			bs, err := d.GetUpBaseInfo(c, mid)
			ctx.Convey("Then err should be nil.bs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(bs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateUpInfo(t *testing.T) {
	convey.Convey("UpdateUpInfo", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(1,2,3,4,5)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpdateUpInfo(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMIDs(t *testing.T) {
	convey.Convey("MIDs", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			offset = int64(0)
			limit  = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			last, mids, err := d.MIDs(c, offset, limit)
			ctx.Convey("Then err should be nil.last,mids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mids, convey.ShouldNotBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
			})
		})
	})
}
