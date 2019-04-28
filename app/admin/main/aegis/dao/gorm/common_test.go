package gorm

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/aegis/model/net"
)

func TestDaoAvailable(t *testing.T) {
	var (
		db = &gorm.DB{}
	)
	convey.Convey("Available", t, func(ctx convey.C) {
		p1 := Available(db)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDisable(t *testing.T) {
	var (
		db = &gorm.DB{}
	)
	convey.Convey("Disable", t, func(ctx convey.C) {
		p1 := Disable(db)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaostatePager(t *testing.T) {
	var (
		s1 = ""
	)
	convey.Convey("state", t, func(ctx convey.C) {
		p1 := state(s1)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDao_ColumnMapString(t *testing.T) {
	var (
		table  = net.TableFlow
		column = "ch_name"
		ids    = []int64{1, 2, 3}
	)

	convey.Convey("ColumnMapString", t, func(ctx convey.C) {
		result, err := d.ColumnMapString(cntx, table, column, ids, "")
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		t.Logf("result(%+v)", result)
	})
}
