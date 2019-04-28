package dao

import (
	"context"
	"reflect"
	"testing"

	"github.com/bouk/monkey"

	"go-common/app/service/main/member/model"
	"go-common/library/cache/memcache"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoexpKey(t *testing.T) {
	var (
		mid = int64(111001740)
	)
	convey.Convey("expKey", t, func(ctx convey.C) {
		p1 := expKey(mid)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaomoralKey(t *testing.T) {
	var (
		mid = int64(111001740)
	)
	convey.Convey("moralKey", t, func(ctx convey.C) {
		p1 := moralKey(mid)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaopingMC(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("pingMC", t, func(ctx convey.C) {
		err := d.pingMC(c)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetBaseInfoCache(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(19476037)
		info = &model.BaseInfo{
			Mid:  19476037,
			Name: "lala",
			Sign: "We are the world!",
		}
	)
	convey.Convey("SetBaseInfoCache", t, func(ctx convey.C) {
		err := d.SetBaseInfoCache(c, mid, info)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoBaseInfoCache(t *testing.T) {
	convey.Convey("BaseInfoCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(19476037)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			info, err := d.BaseInfoCache(c, mid)
			ctx.Convey("Error should be nil. info should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(info, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When conn.Get gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool,
				_ context.Context) memcache.Conn {
				return memcache.MockWith(memcache.ErrItemObject)
			})
			defer guard.Unpatch()
			_, err := d.BaseInfoCache(c, mid)
			ctx.Convey("Error should be equal to memcache.ErrItemObject", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, memcache.ErrItemObject)
			})
		})
	})
}

func TestDaoSetBatchBaseInfoCache(t *testing.T) {
	var (
		c   = context.Background()
		bi1 = &model.BaseInfo{
			Mid:  19476037,
			Name: "lala",
			Sign: "We are the world!",
		}
		bi2 = &model.BaseInfo{
			Mid:  4780461,
			Name: "lala",
			Sign: "We are the world!",
		}
		bs = []*model.BaseInfo{bi1, bi2}
	)
	convey.Convey("SetBatchBaseInfoCache", t, func(ctx convey.C) {
		err := d.SetBatchBaseInfoCache(c, bs)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoBatchBaseInfoCache(t *testing.T) {
	var (
		c    = context.Background()
		mids = []int64{19476037, 1}
	)
	convey.Convey("BatchBaseInfoCache", t, func(ctx convey.C) {
		cached, missed, err := d.BatchBaseInfoCache(c, mids)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("missed should not be nil", func(ctx convey.C) {
			ctx.So(missed, convey.ShouldNotBeNil)
		})
		ctx.Convey("cached should not be nil", func(ctx convey.C) {
			ctx.So(cached, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelBaseInfoCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(19476037)
	)
	convey.Convey("DelBaseInfoCache", t, func(ctx convey.C) {
		err := d.DelBaseInfoCache(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetExpCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(19476037)
		exp = int64(10)
	)
	convey.Convey("SetExpCache", t, func(ctx convey.C) {
		err := d.SetExpCache(c, mid, exp)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoexpCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(19476037)
	)
	convey.Convey("expCache", t, func(ctx convey.C) {
		exp, err := d.expCache(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("exp should not be nil", func(ctx convey.C) {
			ctx.So(exp, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoexpsCache(t *testing.T) {
	var (
		c    = context.Background()
		mids = []int64{19476037, 4780461}
	)
	convey.Convey("expsCache", t, func(ctx convey.C) {
		exps, miss, err := d.expsCache(c, mids)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("miss should not be nil", func(ctx convey.C) {
			ctx.So(miss, convey.ShouldBeNil)
		})
		ctx.Convey("exps should not be nil", func(ctx convey.C) {
			ctx.So(exps, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetMoralCache(t *testing.T) {
	var (
		c     = context.Background()
		mid   = int64(19476037)
		moral = &model.Moral{}
	)
	convey.Convey("SetMoralCache", t, func(ctx convey.C) {
		err := d.SetMoralCache(c, mid, moral)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaomoralCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(19476037)
	)
	convey.Convey("moralCache", t, func(ctx convey.C) {
		moral, err := d.moralCache(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("moral should not be nil", func(ctx convey.C) {
			ctx.So(moral, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelMoralCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(19476037)
	)
	convey.Convey("DelMoralCache", t, func(ctx convey.C) {
		err := d.DelMoralCache(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaorealnameInfoKey(t *testing.T) {
	var (
		mid = int64(19476037)
	)
	convey.Convey("realnameApplyKey", t, func(ctx convey.C) {
		key := realnameInfoKey(mid)
		ctx.Convey("p1 should equal realname_info_<mid>", func(ctx convey.C) {
			ctx.So(key, convey.ShouldEqual, "realname_info_19476037")
		})
	})
}

func TestDaorealnameCaptureTimesKey(t *testing.T) {
	var (
		mid = int64(19476037)
	)
	convey.Convey("realnameCaptureTimesKey", t, func(ctx convey.C) {
		p1 := realnameCaptureTimesKey(mid)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaorealnameCaptureCodeKey(t *testing.T) {
	var (
		mid = int64(19476037)
	)
	convey.Convey("realnameCaptureCodeKey", t, func(ctx convey.C) {
		p1 := realnameCaptureCodeKey(mid)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaorealnameCaptureErrTimesKey(t *testing.T) {
	var (
		mid = int64(19476037)
	)
	convey.Convey("realnameCaptureErrTimesKey", t, func(ctx convey.C) {
		p1 := realnameCaptureErrTimesKey(mid)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetRealnameCaptureTimes(t *testing.T) {
	var (
		c     = context.Background()
		mid   = int64(19476037)
		times = int(0)
	)
	convey.Convey("SetRealnameCaptureTimes", t, func(ctx convey.C) {
		err := d.SetRealnameCaptureTimes(c, mid, times)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRealnameCaptureTimesCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(19476037)
	)
	convey.Convey("RealnameCaptureTimesCache", t, func(ctx convey.C) {
		times, err := d.RealnameCaptureTimesCache(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("times should not be nil", func(ctx convey.C) {
			ctx.So(times, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoIncreaseRealnameCaptureTimes(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(19476037)
	)
	convey.Convey("IncreaseRealnameCaptureTimes", t, func(ctx convey.C) {
		err := d.IncreaseRealnameCaptureTimes(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRealnameCaptureCodeCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(19476037)
	)
	convey.Convey("RealnameCaptureCodeCache", t, func(ctx convey.C) {
		code, err := d.RealnameCaptureCodeCache(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("code should not be nil", func(ctx convey.C) {
			ctx.So(code, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetRealnameInfo(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(19476037)
		info = &model.RealnameCacheInfo{}
	)
	convey.Convey("SetRealnameApplyInfo", t, func(ctx convey.C) {
		err := d.SetRealnameInfo(c, mid, info)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRealnameInfoCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(19476037)
	)
	convey.Convey("RealnameApplyInfoCache", t, func(ctx convey.C) {
		info, err := d.RealnameInfoCache(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("info should not be nil", func(ctx convey.C) {
			ctx.So(info, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDeleteRealnameInfo(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(19476037)
	)
	convey.Convey("DeleteRealnameApplyInfo", t, func(ctx convey.C) {
		err := d.DeleteRealnameInfo(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetRealnameCaptureCode(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(19476037)
		code = int(0)
	)
	convey.Convey("SetRealnameCaptureCode", t, func(ctx convey.C) {
		err := d.SetRealnameCaptureCode(c, mid, code)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDeleteRealnameCaptureCode(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(19476037)
	)
	convey.Convey("DeleteRealnameCaptureCode", t, func(ctx convey.C) {
		err := d.DeleteRealnameCaptureCode(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetRealnameCaptureErrTimes(t *testing.T) {
	var (
		c     = context.Background()
		mid   = int64(19476037)
		times = int(0)
	)
	convey.Convey("SetRealnameCaptureErrTimes", t, func(ctx convey.C) {
		err := d.SetRealnameCaptureErrTimes(c, mid, times)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRealnameCaptureErrTimesCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(19476037)
	)
	convey.Convey("RealnameCaptureErrTimesCache", t, func(ctx convey.C) {
		times, err := d.RealnameCaptureErrTimesCache(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("times should not be nil", func(ctx convey.C) {
			ctx.So(times, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoIncreaseRealnameCaptureErrTimes(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(19476037)
	)
	convey.Convey("IncreaseRealnameCaptureErrTimes", t, func(ctx convey.C) {
		err := d.IncreaseRealnameCaptureErrTimes(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDeleteRealnameCaptureErrTimes(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(19476037)
	)
	convey.Convey("DeleteRealnameCaptureErrTimes", t, func(ctx convey.C) {
		err := d.DeleteRealnameCaptureErrTimes(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
