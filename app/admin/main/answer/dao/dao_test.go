package dao

import (
	"flag"
	"os"
	"testing"

	"go-common/app/admin/main/answer/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account-law.answer-admin")
		flag.Set("conf_appid", "main.account-law.answer-admin")
		flag.Set("conf_token", "bec0ecd7a2799a424602f9a0daea070d")
		flag.Set("tree_id", "4752")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_env", "10")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/answer-admin-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	Convey("Close", t, func() {
		New(conf.Conf)
	})
}

func TestClose(t *testing.T) {
	Convey("Close", t, func() {
		d.Close()
		d = New(conf.Conf)
	})
}
