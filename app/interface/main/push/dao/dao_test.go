package dao

import (
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/interface/main/push/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func init() {
	dir, _ := filepath.Abs("../cmd/push-interface-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func WithDao(f func(d *Dao)) func() {
	return func() {
		f(d)
	}
}

func Test_MiInvalidTokens(t *testing.T) {
	Convey("fetch mi invalid tokens", t, WithDao(func(d *Dao) {
		// 用的时候打开，消息消费完了就没了
		// err := d.DelInvalidMiReports(context.TODO())
		// So(err, ShouldBeNil)
	}))
}
