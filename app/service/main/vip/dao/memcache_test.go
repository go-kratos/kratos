package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/vip/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaopointTip(t *testing.T) {
	var (
		mid = int64(0)
		id  = int64(0)
	)
	convey.Convey("pointTip", t, func(ctx convey.C) {
		p1 := pointTip(mid, id)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoopenCode(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("openCode", t, func(ctx convey.C) {
		p1 := openCode(mid)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaovipInfoKey(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("vipInfoKey", t, func(ctx convey.C) {
		p1 := vipInfoKey(mid)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaosignVip(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("signVip", t, func(ctx convey.C) {
		p1 := signVip(mid)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelVipInfoCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
	)
	convey.Convey("DelVipInfoCache", t, func(ctx convey.C) {
		err := d.DelVipInfoCache(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaopingMC(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("pingMC", t, func(ctx convey.C) {
		err := d.pingMC(c)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetVipInfoCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		v   = &model.VipInfo{}
	)
	convey.Convey("SetVipInfoCache", t, func(ctx convey.C) {
		err := d.SetVipInfoCache(c, mid, v)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoVipInfoCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
	)
	convey.Convey("VipInfoCache", t, func(ctx convey.C) {
		_, err := d.VipInfoCache(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaodelCache(t *testing.T) {
	var (
		c   = context.TODO()
		key = "key"
	)
	convey.Convey("delCache", t, func(ctx convey.C) {
		err := d.DelCache(c, key)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoGetOpenCodeCount(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
	)
	convey.Convey("GetOpenCodeCount", t, func(ctx convey.C) {
		val, err := d.GetOpenCodeCount(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("val should not be nil", func(ctx convey.C) {
			ctx.So(val, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetOpenCode(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(0)
		count = int(0)
	)
	convey.Convey("SetOpenCode", t, func(ctx convey.C) {
		err := d.SetOpenCode(c, mid, count)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoGetPointTip(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		id  = int64(0)
	)
	convey.Convey("GetPointTip", t, func(ctx convey.C) {
		val, err := d.GetPointTip(c, mid, id)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("val should not be nil", func(ctx convey.C) {
			ctx.So(val, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetPointTip(t *testing.T) {
	var (
		c       = context.TODO()
		mid     = int64(0)
		id      = int64(0)
		val     = int(0)
		expired = int32(0)
	)
	convey.Convey("SetPointTip", t, func(ctx convey.C) {
		err := d.SetPointTip(c, mid, id, val, expired)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetSignVip(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		no  = int(0)
	)
	convey.Convey("SetSignVip", t, func(ctx convey.C) {
		err := d.SetSignVip(c, mid, no)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoGetSignVip(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
	)
	convey.Convey("GetSignVip", t, func(ctx convey.C) {
		val, err := d.GetSignVip(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("val should not be nil", func(ctx convey.C) {
			ctx.So(val, convey.ShouldNotBeNil)
		})
	})
}
