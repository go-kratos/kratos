package account

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/up/conf"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/up-admin.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() {})
		f(d)
	}
}

func Test_IdentifyInfo(t *testing.T) {
	var (
		c   = context.Background()
		err error
		mid = int64(27515256)
		ak  = "efbf2e093c3c04008acbd9906ba970db"
		ck  = ""
		ip  = "127.0.0.1"
	)
	Convey("IdentifyInfo", t, WithDao(func(d *Dao) {
		err = d.IdentifyInfo(c, ak, ck, ip, mid)
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ecode.CreativeAccServiceErr)
	}))
}

func Test_GetCachedInfos(t *testing.T) {
	var (
		c    = context.Background()
		err  error
		mids = []int64{1234, 12345}
		ip   = "127.0.0.1"
	)

	Convey("GetCachedInfos", t, WithDao(func(d *Dao) {

		_, err = d.GetCachedInfos(c, mids, ip)
		So(err, ShouldNotBeNil)
	}))
}
