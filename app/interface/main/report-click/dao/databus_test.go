package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPlay(t *testing.T) {
	convey.Convey("Play", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			plat      = ""
			aid       = ""
			cid       = ""
			part      = ""
			mid       = ""
			level     = ""
			ftime     = ""
			stime     = ""
			did       = ""
			ip        = ""
			agent     = ""
			buvid     = ""
			cookieSid = ""
			refer     = ""
			typeID    = ""
			subType   = ""
			sid       = ""
			epid      = ""
			playMode  = ""
			platform  = ""
			device    = ""
			mobiAapp  = ""
			autoPlay  = ""
			session   = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ctx.Convey("No return values", func(ctx convey.C) {
				d.Play(c, plat, aid, cid, part, mid, level, ftime, stime, did, ip, agent, buvid, cookieSid, refer, typeID, subType, sid, epid, playMode, platform, device, mobiAapp, autoPlay, session)
				ctx.SkipSo(" 丹丹让这么写的")
			})
		})
	})
}

func TestDaopubproc(t *testing.T) {
	convey.Convey("pubproc", t, func(ctx convey.C) {
		ctx.Convey("No return values", func(ctx convey.C) {
			d.pubproc()
			ctx.SkipSo(" 丹丹让这么写的")
		})
	})
}

func TestDaomergePub(t *testing.T) {
	var (
		ms = [][]byte{}
	)
	for i := 0; i < 100; i++ {
		var tt = []byte{}
		for j := 0; j < 100; j++ {
			tt = append(tt, []byte("1")...)
		}
		ms = append(ms, tt)
	}
	convey.Convey("mergePub", t, func(ctx convey.C) {
		ctx.Convey("No return values", func(ctx convey.C) {
			d.mergePub(ms)
			ctx.SkipSo(" 丹丹让这么写的")
		})
	})
}
