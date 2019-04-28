package dao

import (
	"context"
	"go-common/app/admin/main/tag/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoFilter(t *testing.T) {
	var (
		c   = context.TODO()
		msg = ""
	)
	convey.Convey("Filter", t, func(ctx convey.C) {
		err := d.Filter(c, msg)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUpSearchTag(t *testing.T) {
	var (
		c   = context.TODO()
		tag = &model.Tag{
			ID:     2,
			State:  -1,
			Verify: 2,
		}
	)
	convey.Convey("UpSearchTag", t, func(ctx convey.C) {
		err := d.UpdateESearchTag(c, tag)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRegionHot(t *testing.T) {
	var (
		c   = context.TODO()
		rid = int64(0)
	)
	convey.Convey("RegionHot", t, func(ctx convey.C) {
		p1, p2, err := d.RegionHot(c, rid)
		ctx.Convey("Then err should be nil.p1,p2 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(p1, convey.ShouldHaveLength, 0)
			ctx.So(p2, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestDaoArchiveHot(t *testing.T) {
	var (
		c   = context.TODO()
		rid = int64(0)
	)
	convey.Convey("ArchiveHot", t, func(ctx convey.C) {
		checked, tags, err := d.ArchiveHot(c, rid)
		ctx.Convey("Then err should be nil.checked,tags should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(tags, convey.ShouldHaveLength, 0)
			ctx.So(checked, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoSendMsg(t *testing.T) {
	var (
		c        = context.TODO()
		mc       = ""
		title    = ""
		context  = ""
		dataType = int32(0)
		mids     = []int64{}
	)
	convey.Convey("SendMsg", t, func(ctx convey.C) {
		err := d.SendMsg(c, mc, title, context, dataType, mids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoBlockUser(t *testing.T) {
	var (
		c               = context.TODO()
		tname           = ""
		title           = ""
		uname           = ""
		action          = ""
		note            = ""
		mid             = int64(0)
		oid             = int64(0)
		reasonType      = int32(0)
		isNotify        = int32(0)
		blockTimeLength = int32(0)
		blockForever    = int32(0)
	)
	convey.Convey("BlockUser", t, func(ctx convey.C) {
		err := d.BlockUser(c, tname, title, uname, action, note, mid, oid, reasonType, isNotify, blockTimeLength, blockForever)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddMoral(t *testing.T) {
	var (
		c        = context.TODO()
		username = ""
		remark   = ""
		reason   = ""
		addMoral = int32(0)
		isNotify = int32(0)
		mids     = []int64{0}
	)
	convey.Convey("AddMoral", t, func(ctx convey.C) {
		err := d.AddMoral(c, username, remark, reason, addMoral, isNotify, mids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}
