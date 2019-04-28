package wall

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/app/conf"
	"go-common/app/admin/main/app/model/wall"

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

func TestWalls(t *testing.T) {
	Convey("get wall all", t, func() {
		res, err := d.Walls(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestWallByID(t *testing.T) {
	Convey("select by id", t, func() {
		res, err := d.WallByID(ctx(), 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestInsert(t *testing.T) {
	Convey("insert wall", t, func() {
		a := &wall.Param{
			Title:    "举杯邀明月",
			Name:     "lllxxx",
			Package:  "对影成三人",
			Logo:     "http://bilibili.com",
			Size:     "25",
			Download: "ssssss",
			Remark:   "sdfsf",
			Rank:     3,
			State:    0,
		}
		err := d.Insert(ctx(), a, time.Now())
		So(err, ShouldBeNil)
	})
}

func TestUpdate(t *testing.T) {
	Convey("update wall", t, func() {
		a := &wall.Param{
			ID:       30,
			Name:     "lllxxx",
			Title:    "举杯邀明月",
			Package:  "对影成三人",
			Logo:     "http://bilibili.com",
			Size:     "25",
			Download: "ssssss",
			Remark:   "sdfsf",
			Rank:     3,
			State:    0,
		}
		err := d.Update(ctx(), a, time.Now())
		So(err, ShouldBeNil)
	})
}

func TestUpdateByID(t *testing.T) {
	Convey("update state by id", t, func() {
		a := &wall.Param{
			IDs: "1,2,3",
		}
		err := d.UpdateByID(ctx(), a, time.Now())
		So(err, ShouldBeNil)
	})
}
