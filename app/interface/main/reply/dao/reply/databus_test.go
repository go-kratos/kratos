package reply

import (
	"context"
	"go-common/app/interface/main/reply/model/reply"
	"go-common/library/queue/databus"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyNewDatabusDao(t *testing.T) {
	convey.Convey("NewDatabusDao", t, func(ctx convey.C) {
		var (
			c = &databus.Config{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			dao := NewDatabusDao(c)
			ctx.Convey("Then dao should not be nil.", func(ctx convey.C) {
				ctx.So(dao, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplypush(t *testing.T) {
	convey.Convey("push", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			key   = ""
			value = interface{}(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Databus.push(c, key, value)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyRecoverFixDialogIdx(t *testing.T) {
	convey.Convey("RecoverFixDialogIdx", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			tp   = int8(0)
			root = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.RecoverFixDialogIdx(c, oid, tp, root)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyRecoverDialogIdx(t *testing.T) {
	convey.Convey("RecoverDialogIdx", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			oid    = int64(0)
			tp     = int8(0)
			root   = int64(0)
			dialog = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.RecoverDialogIdx(c, oid, tp, root, dialog)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyRecoverFloorIdx(t *testing.T) {
	convey.Convey("RecoverFloorIdx", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			oid     = int64(0)
			tp      = int8(0)
			num     = int(0)
			isFloor bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.RecoverFloorIdx(c, oid, tp, num, isFloor)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyDatabusAddTop(t *testing.T) {
	convey.Convey("AddTop", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int8(0)
			top = uint32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.AddTop(c, oid, tp, top)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyDatabusAddReply(t *testing.T) {
	convey.Convey("AddReply", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			rp  = &reply.Reply{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.AddReply(c, oid, rp)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyAddSpam(t *testing.T) {
	convey.Convey("AddSpam", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			mid = int64(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.AddSpam(c, oid, mid, false, tp)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyAddReport(t *testing.T) {
	convey.Convey("AddReport", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			rpID = int64(0)
			tp   = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.AddReport(c, oid, rpID, tp)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyLike(t *testing.T) {
	convey.Convey("Like", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			oid    = int64(0)
			rpID   = int64(0)
			mid    = int64(0)
			action = int8(0)
			ts     = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.Like(c, oid, rpID, mid, action, ts)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyHate(t *testing.T) {
	convey.Convey("Hate", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			oid    = int64(0)
			rpID   = int64(0)
			mid    = int64(0)
			action = int8(0)
			ts     = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.Hate(c, oid, rpID, mid, action, ts)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyRecoverIndex(t *testing.T) {
	convey.Convey("RecoverIndex", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			tp   = int8(0)
			sort = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.RecoverIndex(c, oid, tp, sort)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyRecoverIndexByRoot(t *testing.T) {
	convey.Convey("RecoverIndexByRoot", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			root = int64(0)
			tp   = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.RecoverIndexByRoot(c, oid, root, tp)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyHide(t *testing.T) {
	convey.Convey("Hide", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			rpID = int64(0)
			tp   = int8(0)
			ts   = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.Hide(c, oid, rpID, tp, ts)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyShow(t *testing.T) {
	convey.Convey("Show", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			rpID = int64(0)
			tp   = int8(0)
			ts   = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.Show(c, oid, rpID, tp, ts)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyDelete(t *testing.T) {
	convey.Convey("Delete", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(0)
			oid    = int64(0)
			rpID   = int64(0)
			ts     = int64(0)
			tp     = int8(0)
			assist bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.Delete(c, mid, oid, rpID, ts, tp, assist)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyAdminEdit(t *testing.T) {
	convey.Convey("AdminEdit", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			rpID = int64(0)
			tp   = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.AdminEdit(c, oid, rpID, tp)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyAdminAddTop(t *testing.T) {
	convey.Convey("AdminAddTop", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			adid = int64(0)
			oid  = int64(0)
			rpID = int64(0)
			ts   = int64(0)
			act  = int8(0)
			tp   = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.AdminAddTop(c, adid, oid, rpID, ts, act, tp)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyUpperAddTop(t *testing.T) {
	convey.Convey("UpperAddTop", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(0)
			oid  = int64(0)
			rpID = int64(0)
			ts   = int64(0)
			act  = int8(0)
			tp   = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.UpperAddTop(c, mid, oid, rpID, ts, act, tp)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyAdminDelete(t *testing.T) {
	convey.Convey("AdminDelete", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			adid    = int64(0)
			oid     = int64(0)
			rpID    = int64(0)
			ftime   = int64(0)
			moral   = int(0)
			notify  bool
			adname  = ""
			remark  = ""
			ts      = int64(0)
			tp      = int8(0)
			reason  = int8(0)
			freason = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.AdminDelete(c, adid, oid, rpID, ftime, moral, notify, adname, remark, ts, tp, reason, freason)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyAdminDeleteByReport(t *testing.T) {
	convey.Convey("AdminDeleteByReport", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			adid    = int64(0)
			oid     = int64(0)
			rpID    = int64(0)
			mid     = int64(0)
			ftime   = int64(0)
			moral   = int(0)
			notify  bool
			adname  = ""
			remark  = ""
			ts      = int64(0)
			tp      = int8(0)
			audit   = int8(0)
			reason  = int8(0)
			content = ""
			freason = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.AdminDeleteByReport(c, adid, oid, rpID, mid, ftime, moral, notify, adname, remark, ts, tp, audit, reason, content, freason)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyAdminRecover(t *testing.T) {
	convey.Convey("AdminRecover", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			adid   = int64(0)
			oid    = int64(0)
			rpID   = int64(0)
			remark = ""
			ts     = int64(0)
			tp     = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.AdminRecover(c, adid, oid, rpID, remark, ts, tp)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyAdminPass(t *testing.T) {
	convey.Convey("AdminPass", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			adid   = int64(0)
			oid    = int64(0)
			rpID   = int64(0)
			remark = ""
			ts     = int64(0)
			tp     = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.AdminPass(c, adid, oid, rpID, remark, ts, tp)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyAdminStateSet(t *testing.T) {
	convey.Convey("AdminStateSet", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			adid  = int64(0)
			oid   = int64(0)
			rpID  = int64(0)
			ts    = int64(0)
			tp    = int8(0)
			state = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.AdminStateSet(c, adid, oid, rpID, ts, tp, state)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyAdminTransfer(t *testing.T) {
	convey.Convey("AdminTransfer", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			adid  = int64(0)
			oid   = int64(0)
			rpID  = int64(0)
			ts    = int64(0)
			tp    = int8(0)
			audit = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.AdminTransfer(c, adid, oid, rpID, ts, tp, audit)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyAdminIgnore(t *testing.T) {
	convey.Convey("AdminIgnore", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			adid  = int64(0)
			oid   = int64(0)
			rpID  = int64(0)
			ts    = int64(0)
			tp    = int8(0)
			audit = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.AdminIgnore(c, adid, oid, rpID, ts, tp, audit)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestReplyAdminReportRecover(t *testing.T) {
	convey.Convey("AdminReportRecover", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			adid   = int64(0)
			oid    = int64(0)
			rpID   = int64(0)
			remark = ""
			ts     = int64(0)
			tp     = int8(0)
			audit  = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Databus.AdminReportRecover(c, adid, oid, rpID, remark, ts, tp, audit)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
