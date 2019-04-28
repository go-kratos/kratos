package archive

import (
	"flag"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/videoup/conf"
	"os"
	"testing"
)

var (
	d *Dao
)

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() {})
		f(d)
	}
}
func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.videoup-job")
		flag.Set("conf_token", "c31b4b632890e663002ab6ac22835cda")
		flag.Set("tree_id", "2304")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/videoup-job-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}
