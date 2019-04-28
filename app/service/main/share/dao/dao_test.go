package dao

import (
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/share/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.share-service")
		flag.Set("conf_token", "120ed94e23b963fc0564fbdb662916c3")
		flag.Set("tree_id", "25960")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/share-service-test.toml")
	}

	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func TestPing(t *testing.T) {
	convey.Convey("Ping", t, func() {
		d.Ping()
	})
}
