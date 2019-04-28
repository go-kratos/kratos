package audit

import (
	"context"
	"go-common/app/interface/main/tv/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAuditBeginTrans(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("BeginTran", t, func(ctx convey.C) {
		tx, err := d.BeginTran(c)
		ctx.Convey("Then err should be nil.tx should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tx, convey.ShouldNotBeNil)
		})
		tx.Commit()
	})
}

func TestAuditUpdateVideo(t *testing.T) {
	var (
		c     = context.Background()
		v     = &model.AuditOp{}
		tx, _ = d.BeginTran(c)
	)
	convey.Convey("UpdateVideo", t, func(ctx convey.C) {
		err := d.UpdateVideo(c, v, tx)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		tx.Commit()
	})
}

func TestAuditUpdateArc(t *testing.T) {
	var (
		c     = context.Background()
		v     = &model.AuditOp{}
		tx, _ = d.BeginTran(c)
	)
	convey.Convey("UpdateArc", t, func(ctx convey.C) {
		err := d.UpdateArc(c, v, tx)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		tx.Commit()
	})
}

func TestAuditUpdateCont(t *testing.T) {
	var (
		c     = context.Background()
		val   = &model.AuditOp{}
		tx, _ = d.BeginTran(c)
	)
	convey.Convey("UpdateCont", t, func(ctx convey.C) {
		err := d.UpdateCont(c, val, tx)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		tx.Commit()
	})
}

func TestAuditUpdateSea(t *testing.T) {
	var (
		c     = context.Background()
		val   = &model.AuditOp{}
		tx, _ = d.BeginTran(c)
	)
	convey.Convey("UpdateSea", t, func(ctx convey.C) {
		err := d.UpdateSea(c, val, tx)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		tx.Commit()
	})
}
