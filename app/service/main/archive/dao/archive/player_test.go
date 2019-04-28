package archive

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveValidateQn(t *testing.T) {
	var (
		c  = context.TODO()
		qn = int(80)
	)
	convey.Convey("ValidateQn", t, func(ctx convey.C) {
		allow := d.ValidateQn(c, qn)
		ctx.Convey("Then allow should not be nil.", func(ctx convey.C) {
			ctx.So(allow, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveVipQn(t *testing.T) {
	var (
		c  = context.TODO()
		qn = int(80)
	)
	convey.Convey("VipQn", t, func(ctx convey.C) {
		isVipQn := d.VipQn(c, qn)
		ctx.Convey("Then allow should not be nil.", func(ctx convey.C) {
			ctx.So(isVipQn, convey.ShouldNotBeNil)
		})
	})
}

func TestArchivePlayerInfos(t *testing.T) {
	var (
		c         = context.TODO()
		cids      = []int64{10133958}
		qn        = int(16)
		platform  = "ios"
		ip        = "127.0.0.1"
		fnver     = int(0)
		fnval     = int(8)
		forceHost = int(0)
	)
	convey.Convey("PlayerInfos", t, func(ctx convey.C) {
		d.PlayerInfos(c, cids, qn, platform, ip, fnver, fnval, forceHost)
	})
}

func TestArchivePGCPlayerInfos(t *testing.T) {
	var (
		c        = context.TODO()
		aids     = []int64{10111001}
		platform = "iphone"
		ip       = "121.31.246.238"
		session  = ""
		fnval    = 16
		fnver    = 0
	)
	convey.Convey("PGCPlayerInfos", t, func(ctx convey.C) {
		pm, err := d.PGCPlayerInfos(c, aids, platform, ip, session, fnval, fnver)
		ctx.Convey("Then err should be nil.pm should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(pm, convey.ShouldNotBeNil)
			p, _ := json.Marshal(pm)
			convey.Printf("%s", p)
		})
	})
}
