package bangumi

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/interface/main/web-show/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func WithDao(f func(d *Dao)) func() {
	return func() {
		dir, _ := filepath.Abs("../cmd/web-show-test.toml")
		flag.Set("conf", dir)
		conf.Init()
		if d == nil {
			d = New(conf.Conf)
		}
		f(d)
	}
}

func TestDao_Bangumi(t *testing.T) {
	Convey("test Bangumi", t, WithDao(func(d *Dao) {
		mid := int64(5187977)
		aid := int64(2)
		data := d.IsBp(context.TODO(), mid, aid, "")
		So(data, ShouldNotBeNil)
		Printf("%+v", data)
	}))
}
