package notice

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/app/conf"
	"go-common/app/admin/main/app/model/notice"

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

func TestNotices(t *testing.T) {
	Convey("get Notices all", t, func() {
		res, err := d.Notices(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestNoticeByID(t *testing.T) {
	Convey("get NoticeByID", t, func() {
		res, err := d.NoticeByID(ctx(), 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestInsert(t *testing.T) {
	Convey("insert notice", t, func() {
		a := &notice.Param{
			Plat:    2,
			Title:   "苏轼",
			Content: "不思量，自难忘",
			URL:     "http://www.bilibili.com",
			Area:    "中国台湾",
			Type:    1,
		}
		err := d.Insert(ctx(), a, time.Now())
		So(err, ShouldBeNil)
	})
}

func TestUpdate(t *testing.T) {
	Convey("update notice", t, func() {
		a := &notice.Param{
			ID:      19,
			Plat:    3,
			Title:   "白居易",
			Content: "相见时难别亦难",
			URL:     "http://www.bilibili.com",
			Area:    "中国香港",
			Type:    2,
		}
		err := d.Update(ctx(), a, time.Now())
		So(err, ShouldBeNil)
	})
}

func TestUpdateBuild(t *testing.T) {
	Convey("UpdateBuild notice", t, func() {
		a := &notice.Param{
			ID:         18,
			Build:      2,
			Conditions: "eq",
		}
		err := d.UpdateBuild(ctx(), a, time.Now())
		So(err, ShouldBeNil)
	})
}

func TestUpdateState(t *testing.T) {
	Convey("UpdateRelease notice", t, func() {
		a := &notice.Param{
			ID:    17,
			State: 1,
		}
		err := d.UpdateState(ctx(), a, time.Now())
		So(err, ShouldBeNil)
	})
}
