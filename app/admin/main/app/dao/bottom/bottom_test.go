package bottom

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/app/conf"
	"go-common/app/admin/main/app/model/bottom"

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

func TestBottoms(t *testing.T) {
	Convey("get bottom all", t, func() {
		res, err := d.Bottoms(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestBottomByID(t *testing.T) {
	Convey("selcet bottom by id", t, func() {
		res, err := d.BottomByID(ctx(), 3)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestInsert(t *testing.T) {
	Convey("insert bottom", t, func() {
		a := &bottom.Param{
			Name:   "lllxxx",
			Logo:   "http://i0.hdslb.com/oseA456.jpg",
			Rank:   44,
			Action: 1,
			Param:  "545",
			State:  1,
		}
		err := d.Insert(ctx(), a, time.Now())
		So(err, ShouldBeNil)
	})
}

func TestUpdate(t *testing.T) {
	Convey("update bottom", t, func() {
		a := &bottom.Param{
			ID:     17,
			Name:   "zhanglingyun",
			Logo:   "http://i0.hdslb.com/oseA456.jpg",
			Rank:   44,
			Action: 1,
			Param:  "545",
			State:  1,
		}
		err := d.Update(ctx(), a, time.Now())
		So(err, ShouldBeNil)
	})
}

func TestUpdateByID(t *testing.T) {
	Convey("update by id", t, func() {
		a := &bottom.Param{
			ID:    17,
			State: 0,
		}
		err := d.UpdateByID(ctx(), a, time.Now())
		So(err, ShouldBeNil)
	})
}

func TestDelete(t *testing.T) {
	Convey("delete", t, func() {
		err := d.Delete(ctx(), 15)
		So(err, ShouldBeNil)
	})
}
