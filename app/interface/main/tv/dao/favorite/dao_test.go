package favorite

import (
	"context"
	"encoding/json"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/tv/conf"

	"fmt"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func init() {
	dir, _ := filepath.Abs("../../cmd/tv-interface.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(5 * time.Second)
}

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() {})
		f(d)
	}
}

func TestDao_FavoriteV3(t *testing.T) {
	Convey("TestDao_FavoriteV3", t, func() {
		res, err := d.FavoriteV3(context.Background(), 88894921, 1)
		So(err, ShouldBeNil)
		data, _ := json.Marshal(res)
		Println(string(data))
	})
}

func TestDao_FavAdd(t *testing.T) {
	Convey("TestDao_FavAdd", t, func() {
		err := d.FavAdd(context.Background(), 88894921, 10098813)
		So(err, ShouldBeNil)
		err = d.FavAdd(context.Background(), 88894921, 28417042)
		So(err, ShouldBeNil)
	})
}

func TestDao_FavDel(t *testing.T) {
	Convey("TestDao_FavDel", t, func() {
		err := d.FavDel(context.Background(), 88894921, 28417042)
		So(err, ShouldBeNil)
	})
}

func TestDao_InDefault(t *testing.T) {
	Convey("TestDao_InDefault", t, WithDao(func(d *Dao) {
		var mid = int64(27515418)
		res, err := d.FavoriteV3(context.Background(), mid, 1)
		So(err, ShouldBeNil)
		if res == nil || len(res.List) == 0 {
			fmt.Println("empty Fav")
			return
		}
		exist, err2 := d.InDefault(context.Background(), mid, res.List[0].Oid)
		So(err2, ShouldBeNil)
		So(exist, ShouldBeTrue)
		exist, err2 = d.InDefault(context.Background(), mid, 888888888888)
		So(err2, ShouldBeNil)
		So(exist, ShouldBeFalse)
	}))
}
