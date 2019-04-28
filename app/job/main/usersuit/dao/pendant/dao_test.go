package pendant

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/main/usersuit/conf"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func init() {
	dir, _ := filepath.Abs("../cmd/convey-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func Test_Ping(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.Ping(context.TODO())
		So(err, ShouldBeNil)
	})
}

func Test_Close(t *testing.T) {
	Convey("should return err be nil", t, func() {
		d.Close()
	})
}

func Test_pingRedis(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.pingRedis(context.TODO())
		So(err, ShouldBeNil)
	})
}
