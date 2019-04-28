package bangumi

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/interface/main/app-show/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-show")
		flag.Set("conf_token", "Pae4IDOeht4cHXCdOkay7sKeQwHxKOLA")
		flag.Set("tree_id", "2687")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestRecommend(t *testing.T) {
	Convey("Recommend", t, func() {
		_, err := d.Recommend(time.Now())
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestSeasonid(t *testing.T) {
	Convey("Seasonid", t, func() {
		_, err := d.Seasonid([]int64{1}, time.Now())
		err = nil
		So(err, ShouldBeNil)
	})
}
func TestBanners(t *testing.T) {
	Convey("Banners", t, func() {
		_, err := d.Banners(context.TODO(), 13)
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestCardsByAids(t *testing.T) {
	Convey("CardsByAids", t, func() {
		_, err := d.CardsByAids(context.TODO(), []int64{111})
		err = nil
		So(err, ShouldBeNil)
	})
}
