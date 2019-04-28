package dao

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/main/sms/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func WithDao(f func(d *Dao)) func() {
	return func() {
		dir, _ := filepath.Abs("../cmd/sms-job-test.toml")
		flag.Set("conf", dir)
		flag.Parse()
		conf.Init()
		d = New(conf.Conf)
		f(d)
	}
}

func Test_GetUserMobile(t *testing.T) {
	Convey("get user mobile", t, WithDao(func(d *Dao) {
		mob, err := d.UserMobile(context.TODO(), 27515615)
		So(err, ShouldBeNil)
		So(mob, ShouldNotBeEmpty)
		t.Logf("user(%d) mobile(%v)", 27515615, mob)
	}))
}
