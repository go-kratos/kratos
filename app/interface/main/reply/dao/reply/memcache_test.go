package reply

import (
	"context"
	"go-common/app/interface/main/reply/conf"
	model "go-common/app/interface/main/reply/model/reply"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyNewMemcacheDao(t *testing.T) {
	convey.Convey("NewMemcacheDao", t, func(ctx convey.C) {
		var (
			c = conf.Conf.Memcache
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := NewMemcacheDao(c)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeyCaptcha(t *testing.T) {
	convey.Convey("keyCaptcha", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyCaptcha(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeyAdminTop(t *testing.T) {
	convey.Convey("keyAdminTop", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyAdminTop(oid, tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeyUpperTop(t *testing.T) {
	convey.Convey("keyUpperTop", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyUpperTop(oid, tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeySub(t *testing.T) {
	convey.Convey("keySub", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keySub(oid, tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeyRp(t *testing.T) {
	convey.Convey("keyRp", t, func(ctx convey.C) {
		var (
			rpID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyRp(rpID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplykeyConfig(t *testing.T) {
	convey.Convey("keyConfig", t, func(ctx convey.C) {
		var (
			oid      = int64(0)
			typ      = int8(0)
			category = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyConfig(oid, typ, category)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyMcPing(t *testing.T) {
	convey.Convey("Ping", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Mc.Ping(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyMcCaptchaToken(t *testing.T) {
	convey.Convey("CaptchaToken", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.Mc.CaptchaToken(c, mid)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplySetCaptchaToken(t *testing.T) {
	convey.Convey("SetCaptchaToken", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			token = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Mc.SetCaptchaToken(c, mid, token)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetSubject(t *testing.T) {
	convey.Convey("GetSubject", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			sub, err := d.Mc.GetSubject(c, oid, tp)
			ctx.Convey("Then err should be nil.sub should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(sub, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyGetMultiSubject(t *testing.T) {
	convey.Convey("GetMultiSubject", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oids = []int64{1322313213123}
			tp   = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, missed, err := d.Mc.GetMultiSubject(c, oids, tp)
			ctx.Convey("Then err should be nil.res,missed should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(missed, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyDeleteSubject(t *testing.T) {
	convey.Convey("DeleteSubject", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Mc.DeleteSubject(c, oid, tp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddSubject(t *testing.T) {
	convey.Convey("AddSubject", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			subs = &model.Subject{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Mc.AddSubject(c, subs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddReply(t *testing.T) {
	convey.Convey("AddReply", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			rs = &model.Reply{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Mc.AddReply(c, rs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddTop(t *testing.T) {
	convey.Convey("AddTop", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int8(0)
			rp  = &model.Reply{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Mc.AddTop(c, oid, tp, rp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyDeleteReply(t *testing.T) {
	convey.Convey("DeleteReply", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			rpID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Mc.DeleteReply(c, rpID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetTop(t *testing.T) {
	convey.Convey("GetTop", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int8(0)
			top = uint32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rp, err := d.Mc.GetTop(c, oid, tp, top)
			ctx.Convey("Then err should be nil.rp should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rp, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetReply(t *testing.T) {
	convey.Convey("GetReply", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			rpID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rp, err := d.Mc.GetReply(c, rpID)
			ctx.Convey("Then err should be nil.rp should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rp, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetMultiReply(t *testing.T) {
	convey.Convey("GetMultiReply", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			rpIDs = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rpMap, missed, err := d.Mc.GetMultiReply(c, rpIDs)
			ctx.Convey("Then err should be nil.rpMap,missed should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(missed, convey.ShouldBeNil)
				ctx.So(rpMap, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetReplyConfig(t *testing.T) {
	convey.Convey("GetReplyConfig", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			oid      = int64(0)
			typ      = int8(0)
			category = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			config, err := d.Mc.GetReplyConfig(c, oid, typ, category)
			ctx.Convey("Then err should be nil.config should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(config, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddReplyConfigCache(t *testing.T) {
	convey.Convey("AddReplyConfigCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
			m = &model.Config{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Mc.AddReplyConfigCache(c, m)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
