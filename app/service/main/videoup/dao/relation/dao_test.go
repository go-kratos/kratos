package relation

import (
	"context"
	"flag"
	"go-common/app/service/main/videoup/conf"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/smartystreets/goconvey/convey"
	"path/filepath"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.videoup-service")
		flag.Set("conf_token", "4b62721602981eb3635dba3b0d866ac5")
		flag.Set("tree_id", "2308")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		dir, _ := filepath.Abs("../../cmd/videoup-service.toml")
		flag.Set("conf", dir)
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func TestRelationPing(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		err := d.Ping(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestRelation(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Relation", t, func(ctx convey.C) {
		data, err := d.Relation(c, 27515314, 27515258, 1)
		spew.Dump(data)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAddBlack(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Relation", t, func(ctx convey.C) {
		err := d.AddBalck(c, 27515260, 2089809, 1)
		spew.Dump(err)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
