package dao

import (
	"context"
	"go-common/app/admin/main/reply/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReportLog(t *testing.T) {
	sp := model.LogSearchParam{
		Oid:  10099866,
		Type: 1,
	}
	c := context.Background()
	Convey("test top ReportLog", t, WithDao(func(d *Dao) {
		sp.Action = "top"
		data, err := d.ReportLog(c, sp)
		So(err, ShouldBeNil)
		So(data.Page.Total, ShouldEqual, 1)
	}))
	Convey("test monitor ReportLog", t, WithDao(func(d *Dao) {
		sp.Action = "monitor"
		data, err := d.ReportLog(c, sp)
		So(err, ShouldBeNil)
		So(data.Page.Total, ShouldEqual, 14)
	}))
	Convey("test subject groupby ReportLog", t, WithDao(func(d *Dao) {
		sp.Action = "subject_allow,subject_forbid,subject_frozen,subject_unfrozen_allow,subject_unfrozen_forbid"
		sp.Oid = 0
		data, err := d.ReportLog(c, sp)
		So(err, ShouldBeNil)
		So(data.Page.Total, ShouldEqual, 67)
	}))
	Convey("test one subject ReportLog", t, WithDao(func(d *Dao) {
		sp.Action = "subject_allow,subject_forbid,subject_frozen,subject_unfrozen_allow,subject_unfrozen_forbid"
		sp.Appid = "log_audit"
		sp.Oid = 10099866
		data, err := d.ReportLog(c, sp)
		So(err, ShouldBeNil)
		So(data.Page.Total, ShouldEqual, 2)
	}))
}
