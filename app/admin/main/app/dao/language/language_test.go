package language

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/app/conf"
	"go-common/app/admin/main/app/model/language"

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

func TestLanguages(t *testing.T) {
	Convey("get language all", t, func() {
		res, err := d.Languages(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestLangByID(t *testing.T) {
	Convey("select by id", t, func() {
		res, err := d.LangByID(ctx(), 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestInsert(t *testing.T) {
	Convey("insert language", t, func() {
		a := &language.Param{
			Name:   "4455",
			Remark: "繁体中文",
		}
		err := d.Insert(ctx(), a)
		So(err, ShouldBeNil)
	})
}

func TestUpdate(t *testing.T) {
	Convey("update notice", t, func() {
		a := &language.Param{
			ID:     3,
			Name:   "于谦",
			Remark: "简体中文",
		}
		err := d.Update(ctx(), a)
		So(err, ShouldBeNil)
	})
}
