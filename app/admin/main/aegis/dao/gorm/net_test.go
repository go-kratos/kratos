package gorm

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/aegis/model/net"
)

func TestDaoNetByID(t *testing.T) {
	convey.Convey("NetByID", t, func(ctx convey.C) {
		d.NetByID(cntx, 1)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
		})
	})
}

func TestDao_NetList(t *testing.T) {
	convey.Convey("NetList", t, func(ctx convey.C) {
		pm := &net.ListNetParam{
			BusinessID: 1,
			//State:      net.StateAvailable,
			Ps: 20,
			Pn: 1,
			ID: []int64{1},
		}
		_, err := d.NetList(cntx, pm)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoNetIDByBusiness(t *testing.T) {
	convey.Convey("NetIDByBusiness", t, func(ctx convey.C) {
		res, err := d.NetIDByBusiness(cntx, []int64{1, 2, 3})
		t.Logf("res(%+v)", res)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoNetsByBusiness(t *testing.T) {
	convey.Convey("NetsByBusiness", t, func(ctx convey.C) {
		_, err := d.NetsByBusiness(cntx, []int64{1}, true)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoNets(t *testing.T) {
	convey.Convey("Nets", t, func(ctx convey.C) {
		_, err := d.Nets(cntx, []int64{})
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoNetByUnique(t *testing.T) {
	convey.Convey("NetByUnique", t, func(ctx convey.C) {
		_, err := d.NetByUnique(cntx, "")
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoNetBindStartFlow(t *testing.T) {
	var (
		tx, _ = d.BeginTx(cntx)
	)
	defer tx.Commit()
	convey.Convey("NetBindStartFlow", t, func(ctx convey.C) {
		err := d.NetBindStartFlow(cntx, tx, 0, 0)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDisableNet(t *testing.T) {
	var (
		tx, _ = d.BeginTx(cntx)
	)
	defer tx.Commit()
	convey.Convey("DisableNet", t, func(ctx convey.C) {
		err := d.DisableNet(cntx, tx, 0)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
