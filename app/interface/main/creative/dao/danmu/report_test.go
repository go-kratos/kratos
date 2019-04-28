package danmu

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestDanmuReportUpList(t *testing.T) {
	var (
		c      = context.TODO()
		mid    = int64(2089809)
		pn     = int64(1)
		ps     = int64(10)
		aidStr = "1"
		ip     = "127.0.0.1"
	)
	convey.Convey("ReportUpList", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.dmReportUpListURL).Reply(200).JSON(`{"code":20043,"data":""}`)
		result, total, err := d.ReportUpList(c, mid, pn, ps, aidStr, ip)
		ctx.Convey("Then err should be nil.result,total should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(total, convey.ShouldNotBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestDanmuReportUpArchives(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		ip  = "127.0.0.1"
	)
	convey.Convey("ReportUpArchives", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.dmReportUpArchivesURL).Reply(200).JSON(`{"code":20043,"data":""}`)
		ars, err := d.ReportUpArchives(c, mid, ip)
		ctx.Convey("Then err should be nil.ars should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(ars, convey.ShouldNotBeNil)
		})
	})
}

func TestDanmuReportUpEdit(t *testing.T) {
	var (
		c    = context.TODO()
		mid  = int64(2089809)
		dmid = int64(1)
		cid  = int64(1)
		op   = int64(0)
		ip   = "127.0.0.1"
	)
	convey.Convey("ReportUpEdit", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.dmReportUpArchivesURL).Reply(200).JSON(`{"code":20043,"data":""}`)
		err := d.ReportUpEdit(c, mid, dmid, cid, op, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}
