package upper

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/app-feed/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-feed")
		flag.Set("conf_token", "OC30xxkAOyaH9fI6FRuXA0Ob5HL0f3kc")
		flag.Set("tree_id", "2686")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/app-feed-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
}

func Test_Feed(t *testing.T) {
	Convey("should get Feed", t, func() {
		_, err := d.Feed(context.Background(), 1, 2, 3)
		So(err, ShouldBeNil)
	})
}

func Test_ArchiveFeed(t *testing.T) {
	Convey("should get ArchiveFeed", t, func() {
		_, err := d.ArchiveFeed(context.Background(), 1, 2, 3)
		So(err, ShouldBeNil)
	})
}

func Test_BangumiFeed(t *testing.T) {
	Convey("should get BangumiFeed", t, func() {
		_, err := d.BangumiFeed(context.Background(), 1, 2, 3)
		So(err, ShouldBeNil)
	})
}

func Test_Recent(t *testing.T) {
	Convey("should get Recent", t, func() {
		_, err := d.Recent(context.Background(), 1, 2)
		So(err, ShouldBeNil)
	})
}

func Test_ArticleFeed(t *testing.T) {
	Convey("should get ArticleFeed", t, func() {
		_, err := d.ArticleFeed(context.Background(), 1, 2, 3)
		So(err, ShouldBeNil)
	})
}

func Test_ArticleUnreadCount(t *testing.T) {
	Convey("should get ArticleUnreadCount", t, func() {
		_, err := d.ArticleUnreadCount(context.Background(), 1)
		So(err, ShouldBeNil)
	})
}
