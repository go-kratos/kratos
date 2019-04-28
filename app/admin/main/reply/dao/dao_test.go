package dao

import (
	"flag"
	"os"
	"testing"

	"go-common/app/admin/main/reply/conf"
	"go-common/library/conf/env"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	d  *Dao
	_d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.community.reply-admin")
		flag.Set("conf_token", "7c141ae4ff3f31aade1c51556fd11e8a")
		flag.Set("tree_id", "2124")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		env.DeployEnv = "uat"
		env.Zone = "sh001"
		flag.Set("conf", "../cmd/reply-admin-test.toml")
	}

	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	_d = New(conf.Conf)
	d = _d
	m.Run()
	os.Exit(0)
}

func CleanCache() {

}

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() { CleanCache() })
		f(_d)
	}
}
