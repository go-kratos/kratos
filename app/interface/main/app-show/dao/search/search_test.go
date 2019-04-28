package search

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

func ctx() context.Context {
	return context.Background()
}

func TestMain(m *testing.M) {
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
	} else {
		flag.Set("conf", "../cmd/convey-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	//d.bfsClient.SetTransport(gock.DefaultTransport)
	m.Run()
	os.Exit(0)
}

func TestSearchList(t *testing.T) {
	Convey("SearchList", t, func() {
		_, err := d.SearchList(ctx(), 1, 1, 1, 1, 1, time.Now(), "127.0.0.1", "", "", "", "", "")
		So(err, ShouldBeNil)
	})
}
