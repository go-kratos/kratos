package manager

import (
	"context"
	"go-common/app/service/main/up/dao"
	"go-common/app/service/main/up/model"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestManagerUpSpecials(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("UpSpecials", t, func(ctx convey.C) {
		ups, err := d.UpSpecials(c)
		ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ups, convey.ShouldNotBeNil)
		})
	})
}

func TestManagerRefreshUpSpecialIncremental(t *testing.T) {
	var (
		c         = context.TODO()
		lastMTime = time.Now()
	)
	convey.Convey("RefreshUpSpecialIncremental", t, func(ctx convey.C) {
		ups, err := d.RefreshUpSpecialIncremental(c, lastMTime)
		ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(len(ups), convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestManagerDelSpecialByID(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("DelSpecialByID", t, func(ctx convey.C) {
		res, err := d.DelSpecialByID(c, id)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestManagerInsertSpecial(t *testing.T) {
	var (
		c       = context.TODO()
		special = &model.UpSpecial{}
		mids    = int64(0)
	)
	convey.Convey("InsertSpecial", t, func(ctx convey.C) {
		res, err := d.InsertSpecial(c, special, mids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestManagerUpdateSpecialByID(t *testing.T) {
	var (
		c       = context.TODO()
		id      = int64(0)
		special = &model.UpSpecial{}
	)
	convey.Convey("UpdateSpecialByID", t, func(ctx convey.C) {
		res, err := d.UpdateSpecialByID(c, id, special)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestManagerGetSpecialByMidGroup(t *testing.T) {
	var (
		c       = context.TODO()
		mid     = int64(0)
		groupID = int64(0)
	)
	convey.Convey("GetSpecialByMidGroup", t, func(ctx convey.C) {
		res, err := d.GetSpecialByMidGroup(c, mid, groupID)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestManagerGetSpecialByID(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("GetSpecialByID", t, func(ctx convey.C) {
		res, err := d.GetSpecialByID(c, id)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestManagerGetSepcialCount(t *testing.T) {
	var (
		c          = context.TODO()
		conditions dao.Condition
	)
	convey.Convey("GetSepcialCount", t, func(ctx convey.C) {
		count, err := d.GetSepcialCount(c, conditions)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestManagerGetSpecial(t *testing.T) {
	var (
		c          = context.TODO()
		conditions dao.Condition
	)
	convey.Convey("GetSpecial", t, func(ctx convey.C) {
		res, err := d.GetSpecial(c, conditions)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestManagerGetSpecialByMid(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
	)
	convey.Convey("GetSpecialByMid", t, func(ctx convey.C) {
		res, err := d.GetSpecialByMid(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestManagerRawUpSpecial(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(27515314)
	)
	convey.Convey("RawUpSpecial", t, func(ctx convey.C) {
		res, err := d.RawUpSpecial(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestManagerRawUpsSpecial(t *testing.T) {
	var (
		c    = context.TODO()
		mids = []int64{27515314}
	)
	convey.Convey("RawUpsSpecial", t, func(ctx convey.C) {
		res, err := d.RawUpsSpecial(c, mids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
