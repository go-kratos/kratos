package sports

import (
	"context"
	"flag"
	"net/url"
	"path/filepath"
	"testing"

	"go-common/app/interface/main/activity/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func WithDao(f func(d *Dao)) func() {
	return func() {
		dir, _ := filepath.Abs("../cmd/activity-test.toml")
		flag.Set("conf", dir)
		conf.Init()
		if d == nil {
			d = New(conf.Conf)
		}
		f(d)
	}
}

func TestDao_Qq(t *testing.T) {
	Convey("test dao Qq", t, WithDao(func(d *Dao) {
		var (
			params url.Values
			route  = "matchUnion/fetchData"
		)
		res, err := d.Qq(context.Background(), params, route)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}
