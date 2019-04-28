package dao

import (
	"testing"
	"time"

	"go-common/app/service/main/usersuit/model"

	"github.com/satori/go.uuid"
	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBegin(t *testing.T) {
	convey.Convey("Begin", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.Begin(c)
			ctx.Convey("Then err should be nil.tx should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tx, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxAddInvite(t *testing.T) {
	convey.Convey("TxAddInvite", t, func(ctx convey.C) {
		var (
			tx, _ = d.Begin(c)
			inv   = &model.Invite{
				Mid:  1,
				IPng: []byte{1, 1, 1, 1},
				Code: uuid.NewV4().String(),
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affected, err := d.TxAddInvite(c, tx, inv)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateInvite(t *testing.T) {
	convey.Convey("UpdateInvite", t, func(ctx convey.C) {
		var (
			imid   = int64(0)
			usedAt = int64(0)
			code   = "2bbf90926c984a53"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affected, err := d.UpdateInvite(c, imid, usedAt, code)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInvite(t *testing.T) {
	convey.Convey("Invite", t, func(ctx convey.C) {
		var (
			code = "2bbf90926c984a53"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Invite(c, code)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInvites(t *testing.T) {
	convey.Convey("Invites", t, func(ctx convey.C) {
		var (
			mid = int64(88888970)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Invites(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCurrentCount(t *testing.T) {
	convey.Convey("CurrentCount", t, func(ctx convey.C) {
		var (
			mid   = int64(0)
			start = time.Now()
			end   = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CurrentCount(c, mid, start, end)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
