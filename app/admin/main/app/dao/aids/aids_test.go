package aids

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/app/conf"
	"go-common/app/admin/main/app/model/aids"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-admin-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestInsert(t *testing.T) {
	Convey("insert aids", t, func() {
		a := &aids.Param{
			Aid: 654,
		}
		err := d.Insert(ctx(), a)
		So(err, ShouldBeNil)
	})
}
