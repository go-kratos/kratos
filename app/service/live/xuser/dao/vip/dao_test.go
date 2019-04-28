package vip

import (
	"context"
	"flag"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/xuser/conf"
	"math/rand"
	"path/filepath"
	"testing"
)

// vip dao and conf
var (
	d *Dao
)

type TestUser struct {
	Uid int64
}

// initd init vip dao
func initd() {
	dir, _ := filepath.Abs("../../cmd/test.toml")
	flag.Set("conf", dir)
	flag.Set("deploy_env", "uat")
	conf.Init()
	d = New(conf.Conf)
}

func initTestUser() *TestUser {
	return &TestUser{
		Uid: int64(rand.Int31()),
	}
}

func (t *TestUser) Reset() {
	d.ClearCache(context.Background(), t.Uid)
	d.deleteVip(context.Background(), t.Uid)
}

func testWithTestUser(f func(u *TestUser)) func() {
	u := initTestUser()
	return func() {
		f(u)
		u.Reset()
	}
}

func TestToInt(t *testing.T) {
	var (
		err error
		out int
	)
	Convey("test toInt", t, func() {
		out, err = toInt(1)
		So(out, ShouldEqual, 1)
		So(err, ShouldBeNil)
	})
	Convey("test toInt", t, func() {
		out, err = toInt("123")
		So(out, ShouldEqual, 123)
		So(err, ShouldBeNil)
	})
	Convey("test toInt", t, func() {
		out, err = toInt("test")
		So(out, ShouldEqual, 0)
		So(err, ShouldNotBeNil)
	})
}
