package audit

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/app/conf"
	"go-common/app/admin/main/app/model/audit"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-admin-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestAudits(t *testing.T) {
	Convey("get Audits all", t, func() {
		res, err := d.Audits(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestAuditByID(t *testing.T) {
	Convey("get AuditByID", t, func() {
		res, err := d.AuditByID(ctx(), 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestAuditExist(t *testing.T) {
	Convey("get AuditExist", t, func() {
		a := &audit.Param{
			Build:   222,
			MobiApp: "iphone",
		}
		res, err := d.AuditExist(ctx(), a)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestInsert(t *testing.T) {
	Convey("insert Audit", t, func() {
		a := &audit.Param{
			Build:   222,
			Remark:  "asdsd",
			MobiApp: "iphone",
		}
		err := d.Insert(ctx(), a, time.Now())
		So(err, ShouldBeNil)
	})
}

func TestUpdate(t *testing.T) {
	Convey("update Audit", t, func() {
		a := &audit.Param{
			ID:      1,
			Build:   222,
			Remark:  "asdsd",
			MobiApp: "iphone",
		}
		err := d.Update(ctx(), a, time.Now())
		So(err, ShouldBeNil)
	})
}

func TestDel(t *testing.T) {
	Convey("del Audit by id", t, func() {
		err := d.Del(ctx(), 19)
		So(err, ShouldBeNil)
	})
}
