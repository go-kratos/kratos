package user

import (
	"context"
	"flag"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/live/web-ucenter/conf"
	"go-common/app/interface/live/web-ucenter/dao"
	"testing"
)

var d *Dao

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	flag.Set("env", "uat")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	dao.InitAPI()
	d = New(conf.Conf)
}

var (
	uid = int64(110000681)
	ctx = context.Background()
)

func TestDao_GetAccountProfile(t *testing.T) {
	Convey("test get account profile", t, func() {
		profile, err := d.GetAccountProfile(ctx, uid)
		So(profile, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestDao_GetWallet(t *testing.T) {
	Convey("test get wallet", t, func() {
		_, _, err := d.GetWallet(ctx, uid, "pc")
		So(err, ShouldBeNil)
	})
}

func TestDao_GetLiveAchieve(t *testing.T) {
	Convey("test get rc achieve", t, func() {
		_, err := d.GetLiveAchieve(ctx, uid)
		So(err, ShouldBeNil)
	})
}

func TestDao_GetLiveExp(t *testing.T) {
	Convey("test get exp", t, func() {
		expInfo, err := d.GetLiveExp(ctx, uid)
		So(expInfo, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestDao_GetLiveVip(t *testing.T) {
	Convey("test get vip", t, func() {
		vipInfo, err := d.GetLiveVip(ctx, uid)
		So(vipInfo, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestDao_GetLiveRank(t *testing.T) {
	Convey("test get rankdb", t, func() {
		_, err := d.GetLiveRank(ctx, uid)
		So(err, ShouldBeNil)
	})
}
