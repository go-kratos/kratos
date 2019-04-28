package dao

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/interface/main/kvo/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func WithDao(f func(d *Dao)) func() {
	return func() {
		dir, _ := filepath.Abs("../cmd/web-interface-test.toml")
		flag.Set("conf", dir)
		conf.Init()
		if d == nil {
			d = New(conf.Conf)
		}
		f(d)
	}
}

func TestDao_Document(t *testing.T) {
	Convey("test document", t, WithDao(func(d *Dao) {
		checkSum := int64(11111)
		data, err := d.Document(context.TODO(), checkSum)
		So(err, ShouldBeNil)
		Printf("%+v", data)
	}))
}

func TestDao_UserConf(t *testing.T) {
	Convey("test document", t, WithDao(func(d *Dao) {
		mid := int64(11111)
		moduleKey := 1
		data, err := d.UserConf(context.TODO(), mid, moduleKey)
		So(err, ShouldBeNil)
		Printf("%+v", data)
	}))
}
