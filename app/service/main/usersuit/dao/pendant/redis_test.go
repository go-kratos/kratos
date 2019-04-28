package pendant

import (
	"fmt"
	"go-common/app/service/main/usersuit/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestPendantkeyEquip(t *testing.T) {
	convey.Convey("keyEquip", t, func(ctx convey.C) {
		var (
			mid = int64(650454)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyEquip(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantencode(t *testing.T) {
	convey.Convey("encode", t, func(ctx convey.C) {
		var (
			mid     = int64(650454)
			pid     = int64(1)
			expires = int64(1535970125)
			tp      = int64(0)
			status  = int32(1)
			isVIP   = int32(1)
			pendant = &model.Pendant{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.encode(mid, pid, expires, tp, status, isVIP, pendant)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				fmt.Println(string(res))
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantdecode(t *testing.T) {
	convey.Convey("decode", t, func(ctx convey.C) {
		var (
			src = []byte("\u007b\u0022\u0069\u0064\u0022\u003a\u0030\u002c\u0022\u006d\u0069\u0064\u0022\u003a\u0036\u0035\u0030\u0034\u0035\u0034\u002c\u0022\u0070\u0069\u0064\u0022\u003a\u0031\u002c\u0022\u0065\u0078\u0070\u0069\u0072\u0065\u0022\u003a\u0031\u0035\u0033\u0035\u0039\u0037\u0030\u0031\u0032\u0035\u002c\u0022\u0074\u0079\u0070\u0065\u0022\u003a\u0030\u002c\u0022\u0073\u0074\u0061\u0074\u0075\u0073\u0022\u003a\u0031\u002c\u0022\u0069\u0073\u0056\u0049\u0050\u0022\u003a\u0031\u002c\u0022\u0070\u0065\u006e\u0064\u0061\u006e\u0074\u0022\u003a\u007b\u0022\u0070\u0069\u0064\u0022\u003a\u0030\u002c\u0022\u006e\u0061\u006d\u0065\u0022\u003a\u0022\u0022\u002c\u0022\u0069\u006d\u0061\u0067\u0065\u0022\u003a\u0022\u0022\u002c\u0022\u0069\u006d\u0061\u0067\u0065\u005f\u006d\u006f\u0064\u0065\u006c\u0022\u003a\u0022\u0022\u002c\u0022\u0073\u0074\u0061\u0074\u0075\u0073\u0022\u003a\u0030\u002c\u0022\u0063\u006f\u0069\u006e\u0022\u003a\u0030\u002c\u0022\u0070\u006f\u0069\u006e\u0074\u0022\u003a\u0030\u002c\u0022\u0062\u0063\u006f\u0069\u006e\u0022\u003a\u0030\u002c\u0022\u0065\u0078\u0070\u0069\u0072\u0065\u0022\u003a\u0030\u002c\u0022\u0067\u0069\u0064\u0022\u003a\u0030\u002c\u0022\u0072\u0061\u006e\u006b\u0022\u003a\u0030\u007d\u007d")
			v   = &model.PendantPackage{Mid: 650454, Pid: 1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.decode(src, v)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestPendantAddPKGCache(t *testing.T) {
	convey.Convey("AddPKGCache", t, func(ctx convey.C) {
		var (
			mid  = int64(650454)
			info = []*model.PendantPackage{}
			pp   = &model.PendantPackage{Mid: mid, Pid: int64(1)}
		)
		info = append(info, pp)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddPKGCache(c, mid, info)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestPendantPKGCache(t *testing.T) {
	convey.Convey("PKGCache", t, func(ctx convey.C) {
		var (
			mid = int64(650454)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			info, err := d.PKGCache(c, mid)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(info, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantDelPKGCache(t *testing.T) {
	convey.Convey("DelPKGCache", t, func(ctx convey.C) {
		var (
			mid = int64(650454)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelPKGCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestPendantAddEquipCache(t *testing.T) {
	convey.Convey("AddEquipCache", t, func(ctx convey.C) {
		var (
			mid  = int64(650454)
			info = &model.PendantEquip{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddEquipCache(c, mid, info)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestPendantequipCache(t *testing.T) {
	convey.Convey("equipCache", t, func(ctx convey.C) {
		var (
			mid = int64(650454)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			info, err := d.equipCache(c, mid)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(info, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantDelEquipCache(t *testing.T) {
	convey.Convey("DelEquipCache", t, func(ctx convey.C) {
		var (
			mid = int64(650454)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelEquipCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestPendantequipsCache(t *testing.T) {
	convey.Convey("equipsCache", t, func(ctx convey.C) {
		var (
			mids = []int64{650454}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, p2, err := d.equipsCache(c, mids)
			ctx.Convey("Then err should be nil.p1,p2 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p2, convey.ShouldNotBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantDelEquipsCache(t *testing.T) {
	convey.Convey("DelEquipsCache", t, func(ctx convey.C) {
		var (
			mids = []int64{650454, 1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelEquipsCache(c, mids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
