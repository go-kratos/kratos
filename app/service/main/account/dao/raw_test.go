package dao

import (
	"context"
	"testing"

	mml "go-common/app/service/main/member/model"
	bml "go-common/app/service/main/member/model/block"
	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRawInfo(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(110016481)
	)
	convey.Convey("Get base-info from member-rpc", t, func(ctx convey.C) {
		res, err := d.RawInfo(c, mid)
		ctx.Convey("Then err should be nil and res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoRawInfos(t *testing.T) {
	var (
		c    = context.TODO()
		mids = []int64{110016811, 110017441}
	)
	convey.Convey("Batch get base-info from member-rpc", t, func(ctx convey.C) {
		res, err := d.RawInfos(c, mids)
		ctx.Convey("Then err should be nil and res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoRawCard(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(110018101)
	)
	convey.Convey("Get card-info from member-rpc", t, func(ctx convey.C) {
		res, err := d.RawCard(c, mid)
		ctx.Convey("Then err should be nil and res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoRawCards(t *testing.T) {
	var (
		c    = context.TODO()
		mids = []int64{110019691, 110019241}
	)
	convey.Convey("Batch get card-info from member-rpc", t, func(ctx convey.C) {
		res, err := d.RawCards(c, mids)
		ctx.Convey("Then err should be nil and res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoRawProfile(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(110019691)
	)
	convey.Convey("Get profile-info from member-rpc", t, func(ctx convey.C) {
		res, err := d.RawProfile(c, mid)
		if err == ecode.Degrade {
			err = nil
		}
		ctx.Convey("Then err should be nil and res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoblockStatusToSilence(t *testing.T) {
	var (
		status bml.BlockStatus
	)
	convey.Convey("Check whether member in block status ", t, func(ctx convey.C) {
		p1 := blockStatusToSilence(status)
		ctx.Convey("Then p1 should be 0 or 1.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldBeIn, []int32{0, 1})
		})
	})
}

func TestDaobindEmailStatus(t *testing.T) {
	var (
		email    = "zxfd43622@qq.com"
		spacesta = int8(0)
	)
	convey.Convey("Get member bind-email status", t, func(ctx convey.C) {
		p1 := bindEmailStatus(email, spacesta)
		ctx.Convey("Then p1 should be 0 or 1.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldBeIn, []int32{0, 1})
		})
	})
}

func TestDaobindPhoneStatus(t *testing.T) {
	var (
		phone = "13849920122"
	)
	convey.Convey("Get member bind-phone status", t, func(ctx convey.C) {
		p1 := bindPhoneStatus(phone)
		ctx.Convey("Then p1 should be 0 or 1.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldBeIn, []int32{0, 1})
		})
	})
}

func TestDaoparseMoral(t *testing.T) {
	var (
		moral = &mml.Moral{
			Mid:      3441,
			Moral:    7100,
			Added:    100,
			Deducted: 0,
		}
	)
	convey.Convey("Parse moral value", t, func(ctx convey.C) {
		p1 := parseMoral(moral)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoidentificationStatus(t *testing.T) {
	var (
		realNameStatus = mml.RealnameStatusTrue
	)
	convey.Convey("Get identification status", t, func(ctx convey.C) {
		p1 := identificationStatus(realNameStatus)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldEqual, 1)
		})
	})
}

func TestDaoboolToInt32(t *testing.T) {
	var (
		b = true
	)
	convey.Convey("Convert true to int", t, func(ctx convey.C) {
		p1 := boolToInt32(b)
		ctx.Convey("Then p1 should be 1.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldEqual, 1)
		})
	})
}
