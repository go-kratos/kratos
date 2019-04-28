package data

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
func TestDao_Data(t *testing.T) {
	Convey("test Data", t, WithDao(func(d *Dao) {
		aid := int64(2)
		data, err := d.Related(context.TODO(), aid, "")
		So(err, ShouldBeNil)
		Printf("%+v", data)
	}))
}
