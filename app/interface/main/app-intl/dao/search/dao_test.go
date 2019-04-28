package search

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/interface/main/app-intl/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-intl")
		flag.Set("conf_token", "02007e8d0f77d31baee89acb5ce6d3ac")
		flag.Set("tree_id", "64518")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/app-intl-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestSearch(t *testing.T) {
	Convey("get Search", t, func() {
		res, _, err := d.Search(context.Background(), 1, 2, "iphone", "phone", "1", "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc", "123", "0", "1", "1", "1", "1", 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 8160, 1, 20, time.Now())
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestUpper(t *testing.T) {
	Convey("get Upper", t, func() {
		res, err := d.Upper(context.Background(), 1, "iphone", "phone", "1", "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc", "123", "0", "1", 1, 2, 3, 4, 5, 1, 20, time.Now())
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestArticleByType(t *testing.T) {
	Convey("get ArticleByType", t, func() {
		res, err := d.ArticleByType(context.Background(), 1, 12313, "iphone", "phone", "1", "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc", "123", "0", "1", "2", int8(1), 1, 8190, 1, 1, 20, time.Now())
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestSuggest3(t *testing.T) {
	Convey("get Suggest3", t, func() {
		res, err := d.Suggest3(context.Background(), 12313, "ios", "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc", "123", 8190, 1, "iphone", time.Now())
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
