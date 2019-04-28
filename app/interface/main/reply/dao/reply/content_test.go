package reply

import (
	"context"
	"go-common/app/interface/main/reply/conf"
	"go-common/library/database/sql"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyNewContentDao(t *testing.T) {
	convey.Convey("NewContentDao", t, func(ctx convey.C) {
		var (
			db      = sql.NewMySQL(conf.Conf.MySQL.Reply)
			dbSlave = sql.NewMySQL(conf.Conf.MySQL.ReplySlave)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			dao := NewContentDao(db, dbSlave)
			recover()
			ctx.Convey("Then dao should not be nil.", func(ctx convey.C) {
				ctx.So(dao, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyContenthit(t *testing.T) {
	convey.Convey("hit", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.Content.hit(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyUpMessage(t *testing.T) {
	convey.Convey("UpMessage", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			rpID = int64(0)
			msg  = ""
			now  = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.Content.UpMessage(c, oid, rpID, msg, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyContentGet(t *testing.T) {
	convey.Convey("Get", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			rpID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rc, err := d.Content.Get(c, oid, rpID)
			ctx.Convey("Then err should be nil.rc should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rc, convey.ShouldBeNil)
			})
		})
	})
}

func TestReplyContentGetByIds(t *testing.T) {
	convey.Convey("GetByIds", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(0)
			rpIds = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rcMap, err := d.Content.GetByIds(c, oid, rpIds)
			ctx.Convey("Then err should be nil.rcMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rcMap, convey.ShouldBeNil)
			})
		})
	})
}
