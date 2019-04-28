package reply

import (
	"context"
	model "go-common/app/job/main/reply/model/reply"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplykeyAdminTop(t *testing.T) {
	convey.Convey("keyAdminTop", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
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
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
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
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
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
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyRp(rpID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyPing1(t *testing.T) {
	convey.Convey("Ping", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Ping(c)
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
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Mc.AddSubject(c, subs)
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
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			sub, err := d.Mc.GetSubject(c, oid, tp)
			ctx.Convey("Then err should be nil.sub should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(sub, convey.ShouldNotBeNil)
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
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Mc.AddReply(c, rs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyGetTop1(t *testing.T) {
	convey.Convey("GetTop", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int8(0)
			top = uint32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rp, err := d.Mc.GetTop(c, oid, tp, top)
			ctx.Convey("Then err should be nil.rp should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rp, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyAddTop(t *testing.T) {
	convey.Convey("AddTop", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			rp = &model.Reply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Mc.AddTop(c, rp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyDeleteTop(t *testing.T) {
	convey.Convey("DeleteTop", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			rp = &model.Reply{}
			tp = uint32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Mc.DeleteTop(c, rp, tp)
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
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Mc.DeleteReply(c, rpID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
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
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rp, err := d.Mc.GetReply(c, rpID)
			ctx.Convey("Then err should be nil.rp should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rp, convey.ShouldBeNil)
			})
		})
	})
}
