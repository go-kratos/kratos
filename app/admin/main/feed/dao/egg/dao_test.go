package egg

import (
	"flag"
	"go-common/app/admin/main/feed/conf"
	"os"
	"testing"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app", "feed-admin")
		flag.Set("app_id", "main.web-svr.feed-admin")
		flag.Set("conf_token", "e0d2b216a460c8f8492473a2e3cdd218")
		flag.Set("tree_id", "45266")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}
