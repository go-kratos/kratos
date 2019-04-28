package app

import (
	"context"
	"go-common/app/interface/main/creative/model/app"
	xsql "go-common/library/database/sql"
	"reflect"
	"testing"

	"database/sql"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestAppPortalConfig(t *testing.T) {
	var (
		c  = context.TODO()
		ty = int(0)
	)
	convey.Convey("Portals", t, func(ctx convey.C) {
		apt, err := d.Portals(c, ty)
		ctx.Convey("Then err should be nil.apt should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(apt, convey.ShouldNotBeNil)
		})
	})
}

func TestAppAddMaterialData(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10110788)
		cid = int64(10134702)
		err error
	)
	data := &app.EditorData{
		CID:  cid,
		AID:  aid,
		Type: 3,
		Data: "999",
	}
	convey.Convey("AddMaterialData", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
			return nil, sql.ErrNoRows
		})
		defer guard.Unpatch()
		err = d.AddMaterialData(c, data)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}
