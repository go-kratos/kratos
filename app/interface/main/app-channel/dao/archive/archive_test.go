package archive

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/app-channel/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-channel")
		flag.Set("conf_token", "a920405f87c5bbcca15f3ffebf169c04")
		flag.Set("tree_id", "7852")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/app-view-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestArchive(t *testing.T) {
	Convey("Archive", t, func(ctx C) {
		_, err := d.Archive(context.Background(), 99999999)
		So(err, ShouldNotBeNil) // should 404
	})
}

func TestArchives(t *testing.T) {
	Convey("Archives", t, func(ctx C) {
		_, err := d.Archives(context.Background(), []int64{1})
		So(err, ShouldBeNil)
	})
}

func TestArchivesWithPlayer(t *testing.T) {
	Convey("ArchivesWithPlayer", t, func(ctx C) {
		_, err := d.ArchivesWithPlayer(context.Background(), []int64{1}, 0, "", 0, 0, 0)
		So(err, ShouldBeNil)
	})
}
