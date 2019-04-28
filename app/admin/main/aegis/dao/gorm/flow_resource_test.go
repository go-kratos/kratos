package gorm

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoFRByFlow(t *testing.T) {
	convey.Convey("FRByFlow", t, func(ctx convey.C) {
		_, err := d.FRByFlow(cntx, []int64{})
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoFRByNetRID(t *testing.T) {
	convey.Convey("FRByNetRID", t, func(ctx convey.C) {
		_, err := d.FRByNetRID(cntx, []int64{}, []int64{}, false)
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoFRByUniques(t *testing.T) {
	convey.Convey("FRByUniques", t, func(ctx convey.C) {
		_, err := d.FRByUniques(cntx, []int64{1}, []int64{1}, true)
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoCancelFlowResource(t *testing.T) {
	tx, _ := d.BeginTx(cntx)
	defer tx.Commit()
	convey.Convey("CancelFlowResource", t, func(ctx convey.C) {
		err := d.CancelFlowResource(cntx, tx, []int64{})
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
