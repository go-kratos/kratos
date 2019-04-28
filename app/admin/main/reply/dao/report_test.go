package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/reply/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDataReport(t *testing.T) {
	var (
		typ   = int32(12)
		oid   = int64(909)
		rpID  = int64(111845851)
		oids  = []int64{909}
		rpIDs = []int64{111845851}
		now   = time.Now()
		c     = context.Background()
	)
	Convey("report info", t, WithDao(func(d *Dao) {
		rpt, err := d.Report(c, oid, rpID)
		So(err, ShouldBeNil)
		So(rpt.RpID, ShouldEqual, rpID)
		rpts, err := d.Reports(c, oids, rpIDs)
		So(err, ShouldBeNil)
		So(len(rpts), ShouldEqual, 1)
		rptm, err := d.ReportByOids(c, typ, oids...)
		So(err, ShouldBeNil)
		So(len(rptm), ShouldNotEqual, 0)
	}))
	Convey("report state update", t, WithDao(func(d *Dao) {
		_, err := d.UpReportsState(c, oids, rpIDs, model.ReportStateDelete, now)
		So(err, ShouldBeNil)
		rpt, err := d.Report(c, oids[0], rpIDs[0])
		So(err, ShouldBeNil)
		So(rpt.State, ShouldEqual, model.ReportStateDelete)
	}))
	Convey("report reason update", t, WithDao(func(d *Dao) {
		_, err := d.UpReportsStateWithReason(c, oids, rpIDs, model.ReportStateDelete, 1, "test", now)
		So(err, ShouldBeNil)
		rpt, err := d.Report(c, oids[0], rpIDs[0])
		So(err, ShouldBeNil)
		So(rpt.State, ShouldEqual, model.ReportStateDelete)
		So(rpt.Reason, ShouldEqual, 1)
		So(rpt.Content, ShouldEqual, "test")
	}))
	Convey("report attr update", t, WithDao(func(d *Dao) {
		_, err := d.UpReportsAttrBit(c, oids, rpIDs, model.ReportAttrTransferred, model.AttrYes, now)
		So(err, ShouldBeNil)
		rpt, err := d.Report(c, oids[0], rpIDs[0])
		So(err, ShouldBeNil)
		So(rpt.AttrVal(model.ReportAttrTransferred), ShouldEqual, model.AttrYes)
	}))
}
