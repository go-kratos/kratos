package show

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/app/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func Test_BeginTran(t *testing.T) {
	Convey("BeginTran", t, func() {
		tx, err := d.BeginTran(context.TODO())
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
	})
}

func Test_PTime(t *testing.T) {
	Convey("PTime", t, func() {
		_, err := d.PTime(context.TODO(), time.Now())
		So(err, ShouldBeNil)
	})
}

func Test_Pub(t *testing.T) {
	Convey("Pub", t, func() {
		tx, err := d.BeginTran(context.TODO())
		So(err, ShouldBeNil)
		err = d.Pub(tx, 127)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
	})
}

func Test_PingDB(t *testing.T) {
	Convey("PingDB", t, func() {
		err := d.PingDB(context.TODO())
		So(err, ShouldBeNil)
	})
}
