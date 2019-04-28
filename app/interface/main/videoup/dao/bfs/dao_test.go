package bfs

import (
	"bytes"
	"context"
	"flag"
	"go-common/app/interface/main/videoup/conf"
	"os"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.videoup")
		flag.Set("conf_token", "9772c9629b00ac09af29a23004795051")
		flag.Set("tree_id", "2306")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/videoup.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}
func Test_Upload(t *testing.T) {
	var (
		c    = context.TODO()
		err  error
		loc  string
		body = []byte{}
	)
	convey.Convey("Upload", t, func(ctx convey.C) {
		loc, err = d.Upload(c, "jpeg", bytes.NewReader(body))
		ctx.Convey("Then err should be nil should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(loc, convey.ShouldNotBeNil)
		})
	})
}
