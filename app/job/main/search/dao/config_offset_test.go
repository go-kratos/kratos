package dao

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/main/search/conf"
	"go-common/app/job/main/search/model"

	. "github.com/smartystreets/goconvey/convey"
)

func WithCO(f func(d *Dao)) func() {
	return func() {
		dir, _ := filepath.Abs("../cmd/goconvey.toml")
		flag.Set("conf", dir)
		flag.Parse()
		conf.Init()
		d := New(conf.Conf)
		f(d)
	}
}

func Test_Offset(t *testing.T) {
	Convey("Test_Offset", t, WithCO(func(d *Dao) {
		var (
			err error
			c   = context.TODO()
		)
		d.Offset(c, "", "")
		So(err, ShouldBeNil)
	}))
}

func Test_UpdateOffset(t *testing.T) {
	Convey("Test_UpdateOffset", t, WithCO(func(d *Dao) {
		var (
			err    error
			c      = context.TODO()
			offset = &model.LoopOffset{}
		)
		d.updateOffset(c, offset, "", "")
		So(err, ShouldBeNil)
	}))
}

func Test_BulkInitOffset(t *testing.T) {
	Convey("Test_BulkInitOffset", t, WithCO(func(d *Dao) {
		var (
			c      = context.TODO()
			err    error
			offset = &model.LoopOffset{}
			attrs  = &model.Attrs{
				Table: &model.AttrTable{},
			}
		)
		attrs.Table.TableFrom = 0
		attrs.Table.TableTo = 0
		err = d.bulkInitOffset(c, offset, attrs, []string{})
		So(err, ShouldBeNil)
	}))
}
