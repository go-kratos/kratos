package archive

import (
	"context"
	"go-common/app/interface/main/creative/model/archive"
	xsql "go-common/library/database/sql"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestArchiveFlows(t *testing.T) {
	var (
		c     = context.TODO()
		err   error
		flows []*archive.Flow
	)
	convey.Convey("Flows", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.simpleArchive).Reply(200).JSON(`{"code":20001}`)
		flows, err = d.Flows(c)
		ctx.Convey("Then err should be nil.flows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(flows, convey.ShouldNotBeNil)
		})
	})
}

func TestArchivePorder(t *testing.T) {
	var (
		err error
		c   = context.TODO()
		aid = int64(10110560)
		res *archive.Porder
	)
	convey.Convey("Porder", t, func(ctx convey.C) {
		ret := &xsql.Row{}
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "QueryRow", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) *xsql.Row {
			return ret
		})
		defer guard.Unpatch()
		res, err = d.Porder(c, aid)
		ctx.Convey("Then err should be nil.pd should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
