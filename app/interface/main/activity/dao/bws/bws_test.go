package bws

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/interface/main/activity/conf"

	"go-common/app/interface/main/activity/model/bws"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func WithDao(f func(d *Dao)) func() {
	return func() {
		dir, _ := filepath.Abs("../../cmd/activity-test.toml")
		flag.Set("conf", dir)
		if err := conf.Init(); err != nil {
			panic(err)
		}
		if d == nil {
			d = New(conf.Conf)
		}
		f(d)
	}
}

func TestDao_Binding(t *testing.T) {
	Convey("test binding", t, WithDao(func(d *Dao) {
		key := "9875fa517967622b"
		err := d.Binding(context.Background(), 908087, &bws.ParamBinding{Key: key})
		So(err, ShouldBeNil)
	}))
}

func TestDao_UsersKey(t *testing.T) {
	Convey("test users key", t, WithDao(func(d *Dao) {
		key := "9875fa517967622b"
		data, err := d.UsersKey(context.Background(), key)
		So(err, ShouldBeNil)
		Printf("%+v", data)
	}))
}

func TestDao_UsersMid(t *testing.T) {
	Convey("test users mid", t, WithDao(func(d *Dao) {
		mid := int64(908087)
		data, err := d.UsersMid(context.Background(), mid)
		So(err, ShouldBeNil)
		Printf("%+v", data)
	}))
}

func TestDao_CacheUsersMid(t *testing.T) {
	Convey("test cache users mid", t, WithDao(func(d *Dao) {
		mid := int64(908087)
		data, err := d.CacheUsersMid(context.Background(), mid)
		So(err, ShouldBeNil)
		Printf("%+v", data)
	}))
}

func TestDao_DelCacheUsersMid(t *testing.T) {
	Convey("test delete users mid", t, WithDao(func(d *Dao) {
		mid := int64(908087)
		err := d.DelCacheUsersMid(context.Background(), mid)
		So(err, ShouldBeNil)
	}))
}

func TestDao_RawUsersKey(t *testing.T) {
	Convey("test users mid", t, WithDao(func(d *Dao) {
		keyID := int64(1)
		data, err := d.UserByID(context.Background(), keyID)
		So(err, ShouldBeNil)
		Printf("%+v", data)
	}))
}
