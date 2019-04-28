package dao

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/admin/main/point/conf"
	"go-common/app/admin/main/point/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../cmd/point-admin.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func Test_PointHistory(t *testing.T) {
	Convey("Test_PointHistory", t, func() {
		arg := &model.ArgPointHistory{}
		res, err := d.PointHistory(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
