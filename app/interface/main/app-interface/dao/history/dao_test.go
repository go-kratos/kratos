package history

import (
	"context"
	"go-common/app/interface/main/app-interface/conf"
	"testing"

	"flag"
	"path/filepath"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func init() {
	dir, _ := filepath.Abs("../cmd/app-interface-test.toml")
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

func TestDao_ArchiveInfo(t *testing.T) {
	Convey("ArchiveInfo", t, WithDao(func(d *Dao) {
		_, err := d.Archive(context.TODO(), []int64{1, 2})
		So(err, ShouldBeNil)
	}))
}

func TestDao_ArticleInfo(t *testing.T) {
	Convey("ArticleInfo", t, WithDao(func(d *Dao) {
		_, err := d.Article(context.TODO(), []int64{1, 2})
		So(err, ShouldBeNil)
	}))
}

func TestDao_PGCInfo(t *testing.T) {
	Convey("PGCInfo", t, WithDao(func(d *Dao) {
		_, err := d.PGC(context.TODO(), "1", "1", 111, 27515256)
		So(err, ShouldBeNil)
	}))
}

func TestDao_GetList(t *testing.T) {
	Convey("GetList", t, WithDao(func(d *Dao) {
		_, err := d.History(context.TODO(), 27515256, 1, 20)
		So(err, ShouldBeNil)
	}))
}
